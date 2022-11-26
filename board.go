package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

const (
	StartingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type Spaces [64]Piece
type PsuedoMoves [64]Bit
type Hash uint64

const (
	zBlack = iota + 12*64
	zWOO
	zWOOO
	zBOO
	zBOOO
	zEP
)

var zLookups [12*64 + 1 + 4 + 8]Hash

// Set sets a position as attacked.
func (a *PsuedoMoves) Update(p Piece, c Coord) {
	bit := c.Bit()

	// If we're clearing a piece, just clear all spaces on the board it attacks.
	if p == Empty {
		for i := 0; i < 64; i++ {
			a[i] &= ^bit
		}
		return
	}

	// Set the attacked bits for the given piece.
	for _, d := range p.AttackDir() {
		for i, dis, pos := 0, p.AttackDistance(d), c; i < dis; i++ {
			pos = pos.ApplyDir(d)
			if !pos.IsValid() {
				break
			}
			a[pos.Idx()] |= bit
		}
	}
}

// PossibleMove returns true from attacks to.
func (a *PsuedoMoves) PossibleMove(from, to Coord) bool {
	return (a[to.Idx()] & from.Bit()) != 0
}

// Attackers returns a slice of Coord for all attacking squares.
func (a *PsuedoMoves) Attackers(c Coord) Bit {
	return a[c.Idx()]
}

func (a PsuedoMoves) String() string {
	var str string
	for i := 0; i < 64; i++ {
		str += fmt.Sprintf("%02d\t%064b\n", i, a[i])
	}
	return str
}

func (m *PsuedoMoves) countZeros() (total int) {
	for _, v := range m {
		if v != 0 {
			total += 1
		}
	}
	return total
}

// BoardState contains the state of the board that's Undoable.
type BoardState struct {
	turn                 Piece // Whose turn is it?
	wOO, wOOO, bOO, bOOO bool  // Can the white/black king castle kingside/queenside?
	wkLoc, bkLoc         Coord // Where are the kings located?
	fullMove, halfMove   int
	epTarget             Coord
	hash                 Hash
	score                Score
}

type Board struct {
	spaces   Spaces
	state    BoardState
	moves    []Move
	oldState []BoardState

	// State of the kings.
	isWCheck, isBCheck bool

	// Bitfields stating if white or black attack a given square.
	wPseudos PsuedoMoves
	bPseudos PsuedoMoves
}

func (b *Board) at(c Coord) Piece {
	return b.spaces[c.Idx()]
}

func (b *Board) set(p Piece, c Coord) {
	idx := c.Idx()

	// Update the score.
	b.state.score -= b.at(c).Score()
	b.state.score += p.Score()

	// Update the hash.
	hashP := p
	if hashP == Empty {
		hashP = b.at(c)
	}
	if hashP != Empty {
		b.state.hash ^= zLookups[hashP.HashIdx()+idx]
	}

	// Update the king's location.
	if p.Colorless() == King {
		if p.Color() == White {
			b.state.wkLoc = c
		} else {
			b.state.bkLoc = c
		}
	}
	if p == Empty {
		b.PseudoMoves(White).Update(p, c)
		b.PseudoMoves(Black).Update(p, c)
	} else {
		b.PseudoMoves(p).Update(p, c)
	}
	b.spaces[idx] = p
}

// ZHash returns the Zobrist hash for this board state.
func (b *Board) ZHash() Hash {
	return b.state.hash
}

// Print pretty-prints a board on standard out.
func (b *Board) Print() {
	for y := 7; y >= 0; y-- {
		println("  +-+-+-+-+-+-+-+-+")
		print(y+1, " ")
		for x := 0; x < 8; x++ {
			print("|", b.at(Coord{x, y}).String())
		}
		println("|")
	}
	println("  +-+-+-+-+-+-+-+-+")
	println("   a b c d e f g h ")
}

func (b *Board) castleString() string {
	var s string
	canCastle := false
	if b.state.wOO {
		canCastle = true
		s += "K"
	}
	if b.state.wOOO {
		canCastle = true
		s += "Q"
	}
	if b.state.bOO {
		canCastle = true
		s += "k"
	}
	if b.state.bOOO {
		canCastle = true
		s += "q"
	}
	if canCastle == false {
		s = "-"
	}
	return s
}

// String returns the FEN string for a board.
func (b *Board) String() string {
	return b.FENString()
}

// FENString returns the FEN string for the board.
// https://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation
func (b *Board) FENString() string {
	var s string
	var empty int
	printEmpty := func() {
		if empty != 0 {
			s += fmt.Sprintf("%d", empty)
		}
		empty = 0
	}

	// Print the board.
	for y := 7; y >= 0; y-- {
		for x := 0; x <= 7; x++ {
			c := Coord{x, y}
			if p := b.at(c); p.IsEmpty() {
				empty += 1
				continue
			} else {
				printEmpty()
				s += fmt.Sprintf("%v", p)
			}
		}
		printEmpty()
		if y != 0 {
			s += "/"
		}
	}
	s += " "

	// Print the turn
	if b.state.turn == White {
		s += "w "
	} else {
		s += "b "
	}

	// Print the Casting rights
	s += b.castleString() + " "

	// Print the rest.
	s += fmt.Sprintf("%v %d %d", b.state.epTarget, b.state.halfMove, b.state.fullMove)
	return s
}

// PseudoMoves returns the psudeo squares for a given color.
func (b *Board) PseudoMoves(p Piece) *PsuedoMoves {
	if p == Empty {
		panic("empty")
	}
	if p.Color() == White {
		return &b.wPseudos
	}
	return &b.bPseudos
}

// IsKingInCheck returns true if the given piece's color's king is in check.
func (b *Board) IsKingInCheck(p Piece) bool {
	if p == Empty {
		panic("empty")
	}
	if p.Color() == White {
		return b.isWCheck
	}
	return b.isBCheck
}

// KingLoc returns the king location for the given piece's color.
func (b *Board) KingLoc(p Piece) Coord {
	if p.Color() == White {
		return b.state.wkLoc
	}
	return b.state.bkLoc
}

// setCheck sets the check state.
func (b *Board) setCheck(p Piece, v bool) {
	if p.Color() == White {
		b.isWCheck = v
	} else {
		b.isBCheck = v
	}
}

// doesSquareAttack returns true if a given square attacks another one.
// if the square is definitely attacked.
func (b *Board) doesSquareAttack(from, to Coord, color Piece) bool {
	p := b.at(from)
	// If there's no piece at the attacking square, clearly no.
	if p == Empty {
		return false
	}

	// Make sure we could get between the two squares.
	// NB: We could check if d is in the set of attacking squares from that the
	// piece has, but the DoesAttack call below does that as well.
	d := DirBetween(from, to)
	if d == InvalidDir {
		return false
	}

	// Now, check the psuedo-moves, and see if it's a possibility that the
	// piece attacks the given square.
	if !b.PseudoMoves(p).PossibleMove(from, to) {
		return false
	}

	// And finally, make sure there's no pieces between us.
	if !d.IsKnight() {
		for last := from.ApplyDir(d); last != to; last = last.ApplyDir(d) {
			if !last.IsValid() {
				panic("shouldn't be invalid")
			}
			if b.at(last) != Empty {
				return false
			}
		}
	}
	return true
}

// isSquareAttacked returns true if a given square is attacked by the given color.
func (b *Board) isSquareAttacked(c Coord, color Piece) bool {
	if c == InvalidCoord {
		return false
	}

	at := b.PseudoMoves(color).Attackers(c)
	for i := 0; i < 64; i++ {
		t := Bit(1) << i
		if at&t != 0 {
			if b.doesSquareAttack(CoordFromBit(t), c, color) {
				return true
			}
		}
	}
	return false
}

// wouldKingBeInCheck returns true if a move would result in an illegal check.
func (b *Board) wouldKingBeInCheck(m *Move) bool {
	// So, if the king was in check, we need to see if we would block the
	// check, or take the checking piece.
	b.MakeMove(*m)
	check := b.IsKingInCheck(m.p)
	if b.IsKingInCheck(m.p.OppositeColor()) {
		m.isCheck = true
	}
	b.UnmakeMove()
	return check
}

// isLegalMove returns true if we're dealing with a legal move. Also, sets the
// capture state on the move if it would be one.
func (b *Board) isLegalMove(m *Move) bool {
	// Any moves outside the board are invalid.
	if !m.to.IsValid() {
		return false
	}

	// Check that the piece is at the place specified
	if p := b.at(m.from); p != m.p {
		panic("piece was not present")
	}

	p2 := b.at(m.to)

	if m.p.IsKing() {
		// Can't move into check
		if b.isSquareAttacked(m.to, m.p.OppositeColor()) {
			return false
		}

		if m.IsCastle() {
			if m.p.Color() == White {
				if m.IsKingsideCastle() && !b.state.wOO {
					return false
				}
				if m.IsQueensideCastle() && !b.state.wOOO {
					return false
				}
			} else {
				if m.IsKingsideCastle() && !b.state.bOO {
					return false
				}
				if m.IsQueensideCastle() && !b.state.bOOO {
					return false
				}
			}

			// Can't castle out of check.
			if b.IsKingInCheck(m.p) {
				return false
			}

			mid := m.CastleMidCoord()

			// Can't castle across or into non-empty squares.
			if p2 != Empty || b.at(mid) != Empty {
				return false
			}
			if m.IsQueensideCastle() { // also check on knight for O-O-O.
				if b.at(Coord{m.to.x - 1, m.to.y}) != Empty {
					return false
				}
			}

			// Can't castle across a square in check.
			if b.isSquareAttacked(mid, m.p.OppositeColor()) {
				return false
			}
		}
	}

	if m.p.IsPawn() {
		if m.IsVertical() {
			dist := m.to.y - m.from.y
			if dist == 2 || dist == -2 {
				// Pawns can't move 2 spaces if it's not from the start location.
				if m.p.Color() == White && m.from.y != 1 ||
					m.p.Color() == Black && m.from.y != 6 {
					return false
				}

				// Pawns can't move through a piece.
				mid := Coord{m.from.x, (m.from.y + m.to.y) / 2}
				if b.at(mid) != Empty {
					return false
				}
			}

			// Vertical moves must be into an empty space.
			if p2 != Empty {
				return false
			}
		} else if p2 == Empty {
			// Otherwise captures that are empty must be en passant
			if m.to != b.state.epTarget {
				return false
			}
			m.isEnPassant = true
			m.captured = Pawn | m.p.OppositeColor()
			m.isCapture = true
		} else if p2.Color() == m.p.Color() {
			// Cannot capture your own color.
			return false
		} else {
			m.captured = p2
			m.isCapture = true
		}
	} else if p2 != Empty {
		// Can't take pieces of your own color.
		if p2.Color() == m.p.Color() {
			return false
		}
		m.captured = p2
		m.isCapture = true
	}

	// And finally, if the move would result in the king being in check, it's illegal.
	if b.wouldKingBeInCheck(m) {
		return false
	}

	return true
}

// GetMoves returns all moves for a given coordinate.
func (b *Board) GetMoves(moves []Move, c Coord) []Move {
	p := b.at(c)
	if p.IsEmpty() {
		return moves
	}

	for _, d := range p.MoveDir() {
		lastPos := c
		for i, dis := 0, p.SlideDistance(); i < dis; i++ {
			// Skip positions outside the board.
			toPos := lastPos.ApplyDir(d)
			if !toPos.IsValid() {
				break
			}

			// If we'd overlap one our own pieces, no need to check further.
			if b.at(toPos).Color() == p.Color() {
				break
			}

			move := Move{
				p:    p,
				to:   toPos,
				from: c,
			}
			lastPos = toPos

			// If the move would be illegal, we keep checking as a different move in
			// this direction might be legal.
			if !b.isLegalMove(&move) {
				if b.at(toPos) == Empty {
					continue
				} else {
					break
				}
			}

			// It's a legal move.
			if !move.IsPromotion() {
				moves = append(moves, move)
			} else {
				// Add a promotion for each piece.
				for _, promotion := range []Piece{Knight, Bishop, Rook, Queen} {
					move.promotion = promotion | p.Color()
					moves = append(moves, move)
				}
			}
			if move.isCapture { // stop looking when you'd take a piece.
				break
			}
		}
	}
	return moves
}

// PossibleMoves returns a slice of the possible moves for a given Board.
func (b *Board) PossibleMoves(moves []Move) []Move {
	for p := 0; p < 64; p++ {
		c := CoordFromIdx(p)
		p := b.at(c)
		if p.Color() != b.state.turn { // not the color we want, skip
			continue
		}
		moves = b.GetMoves(moves, c)
	}
	return moves
}

// epTarget returns the en passant target of a Move if the move was a pawn
// move, otherwise it returns InvalidCoord
func epTarget(m Move) Coord {
	distY := (m.to.y - m.from.y)
	if m.p.IsPawn() && (distY == 2 || distY == -2) {
		return Coord{m.to.x, (m.to.y + m.from.y) / 2}
	}
	return InvalidCoord
}

// handleRookMoveOrCapture updates the castle state if a rook is moved or
// captured.
func (b *Board) handleRookMoveOrCapture(c Coord) {
	p := b.at(c)
	if !p.IsRook() {
		return
	}

	if p.IsWhite() && c.y == 0 {
		if c.x == 0 {
			b.state.wOOO = false
		} else if c.x == 7 {
			b.state.wOO = false
		}
	} else if p.IsBlack() && c.y == 7 {
		if c.x == 0 {
			b.state.bOOO = false
		} else if c.x == 7 {
			b.state.bOO = false
		}
	}
}

// zCastle returns the Zobrist hash for the given castle state.
func (b *Board) zCastle() (v Hash) {
	if b.state.wOO {
		v ^= zWOO
	}
	if b.state.wOOO {
		v ^= zWOOO
	}
	if b.state.bOO {
		v ^= zBOO
	}
	if b.state.bOOO {
		v ^= zBOOO
	}
	return v
}

// updateCastleState updates the castling state for a given move.
func (b *Board) updateCastleState(m Move) {
	b.state.hash ^= b.zCastle()
	// If we're moving a king, we can't castle anymore.
	if m.p.IsKing() {
		if m.p.IsWhite() {
			b.state.wOO, b.state.wOOO = false, false
		} else {
			b.state.bOO, b.state.bOOO = false, false
		}
	}

	// If we're moving or capturing a rook, update the state.
	b.handleRookMoveOrCapture(m.from)
	b.handleRookMoveOrCapture(m.to)
	b.state.hash ^= b.zCastle()
}

// updateEPTarget updates the enpassant target.
func (b *Board) updateEPTarget(m Move) {
	if b.state.epTarget != InvalidCoord {
		b.state.hash ^= zLookups[zEP+b.state.epTarget.FileIdx()]
	}
	b.state.epTarget = epTarget(m)
	if b.state.epTarget != InvalidCoord {
		b.state.hash ^= zLookups[zEP+b.state.epTarget.FileIdx()]
	}
}

// updateChecks updates the checks.
func (b *Board) updateChecks() {
	b.setCheck(White|King, b.isSquareAttacked(b.KingLoc(White|King), Black))
	b.setCheck(Black|King, b.isSquareAttacked(b.KingLoc(Black|King), White))
}

// MakeMove applies the move, and updates all necessary Board state.
func (b *Board) MakeMove(m Move) {
	// Save some state so we can undo the move if asked.
	b.oldState = append(b.oldState, b.state)

	// Update turn variables, and board state.
	b.state.hash ^= zLookups[zBlack]
	if b.state.turn == White {
		b.state.turn = Black
	} else {
		b.state.fullMove += 1
		b.state.turn = White
	}
	b.state.halfMove += 1
	if m.p.IsPawn() || m.isCapture {
		b.state.halfMove = 0
	}
	b.updateCastleState(m)
	b.updateEPTarget(m)

	// Move the piece.
	b.set(Empty, m.from)
	if m.isCapture {
		b.set(Empty, m.to)
	}
	if !m.IsPromotion() {
		if m.IsCastle() {
			dist := m.to.x - m.from.x
			rookFrom, rookTo := Coord{7, m.from.y}, Coord{5, m.from.y}
			if dist < 0 { // queenside
				rookFrom.x, rookTo.x = 0, 3
			}
			b.set(Empty, rookFrom)
			b.set(m.p.Color()|Rook, rookTo)
		}
		if m.isEnPassant { // Need to remove captured pawn.
			c := m.to
			if m.p.Color() == White {
				c.y = 4
			} else {
				c.y = 3
			}
			b.set(Empty, c)
		}
		b.set(m.p, m.to)
	} else {
		b.set(m.promotion, m.to)
	}

	b.updateChecks()

	// Save off the move.
	b.moves = append(b.moves, m)
}

// UnmakeMove undoes the last move.
func (b *Board) UnmakeMove() {
	// Can't undo a move if we don't have any.
	if len(b.moves) == 0 {
		return
	}

	// Pop the last move.
	m := b.moves[len(b.moves)-1]
	b.moves = b.moves[:len(b.moves)-1]

	// And fix board state.
	b.set(Empty, m.to)
	b.set(m.p, m.from)
	if m.isCapture {
		l := m.to
		if m.isEnPassant {
			if m.p.Color() == White {
				l.y -= 1
			} else {
				l.y += 1
			}
		}
		b.set(m.captured, l)
	} else if c := m.RookCoord(); c != InvalidCoord {
		b.set(Empty, m.CastleMidCoord())
		b.set(m.p.Color()|Rook, m.RookCoord())
	}

	b.updateChecks()

	// Fix game state.
	b.state = b.oldState[len(b.oldState)-1]
	b.oldState = b.oldState[:len(b.oldState)-1]
}

// GetMove gets a move given two coordinates.
func (b *Board) GetMove(from, to Coord) (Move, error) {
	if !from.IsValid() {
		return Move{}, errors.New("invalid from coord")
	}
	if !to.IsValid() {
		return Move{}, errors.New("invalid to coord")
	}
	for _, m := range b.GetMoves(nil, from) {
		if m.from != from || m.to != to {
			continue
		}
		return m, nil
	}
	return Move{}, fmt.Errorf("invalid move: %v", Move{from: from, to: to})
}

type perftHash struct {
	h Hash
	d int
}

// Perft calculates the number of possible moves at a given depth. It's quite
// helpful debugging the move generation. Optionally, Perft will also print the
// number of reachable moves for each valid move in the given board state.
func (b *Board) Perft(origDepth int) uint64 {
	if origDepth == 0 {
		return 0
	}

	// Prevents allocating space for moves at every depth.
	moveQueue := make([][]Move, origDepth)

	// Keep around a set of counts.
	counts := make(map[perftHash]uint64)

	var perft func(int, bool) uint64
	perft = func(d int, s bool) uint64 {
		// Exit early.
		moves := moveQueue[origDepth-d][:0]
		moves = b.PossibleMoves(moves)
		if d <= 1 {
			if s {
				for _, m := range moves {
					fmt.Printf("%v: 1\n", m)
				}
			}
			return uint64(len(moves))
		}

		// Haven't seen this position before, need to calculate it.
		var total uint64
		for _, move := range moves {
			b.MakeMove(move)

			var cnt uint64
			h := perftHash{b.state.hash, d}
			if v, ok := counts[h]; ok {
				cnt = v
			} else {
				cnt = perft(d-1, false)
			}
			counts[h] = cnt

			if s {
				fmt.Printf("%v: %d\n", move, cnt)
			}
			b.UnmakeMove()
			total += cnt
		}
		return total
	}

	return perft(origDepth, false)
}

// EmptyBoard returns a new, empty board. No state of gameplay is set up.
func EmptyBoard() *Board {
	return &Board{
		state: BoardState{
			epTarget: InvalidCoord,
			wkLoc:    InvalidCoord,
			bkLoc:    InvalidCoord,
			fullMove: 1,
		},
	}
}

// New returns a new Board, set up for play (ie a new chess game).
func New() *Board {
	b, err := FromFEN(StartingFEN)
	if err != nil {
		panic(err)
	}
	return b
	/*
		b := EmptyBoard()
		b.state.turn = White
		b.state.wOO = true
		b.state.wOOO = true
		b.state.bOO = true
		b.state.bOOO = true
		pieces := []Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
		for x, p := range pieces {
			b.set(p|White, Coord{x, 0})
			b.set(Pawn|White, Coord{x, 1})
			b.set(Pawn|Black, Coord{x, 6})
			b.set(p|Black, Coord{x, 7})
		}
		return b
	*/
}

var runeToPiece = map[rune]Piece{
	'r': Rook | Black,
	'n': Knight | Black,
	'b': Bishop | Black,
	'q': Queen | Black,
	'k': King | Black,
	'p': Pawn | Black,
	'R': Rook | White,
	'N': Knight | White,
	'B': Bishop | White,
	'Q': Queen | White,
	'K': King | White,
	'P': Pawn | White,
}

func (b *Board) validate() error {
	if b.state.wkLoc == InvalidCoord {
		return errors.New("no white king on board")
	}
	if b.state.bkLoc == InvalidCoord {
		return errors.New("no black king on board")
	}
	return nil
}

// FromFEN creates a Board from a FEN string.
func FromFEN(s string) (*Board, error) {
	b := EmptyBoard()
	coord := Coord{0, 7}
	parts := strings.Fields(s)

	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid number of FEN fields: %d", len(parts))
	}

	// Parse board.
	for _, c := range parts[0] {
		var inc int
		if unicode.IsNumber(c) {
			inc = int(c - '0')
		} else if c != '/' {
			inc = 1
			if p, ok := runeToPiece[c]; ok {
				b.set(p, coord)
			} else {
				return nil, errors.New(fmt.Sprintf("didn't find %v", c))
			}
		}
		coord = Coord{
			(coord.x + inc) % 8,
			(coord.y - (coord.x+inc)/8),
		}
	}

	// Parse turn.
	b.state.turn = White
	if parts[1] == "b" {
		b.state.turn = Black
	}

	// Parse Castling.
	for _, c := range parts[2] {
		switch c {
		case 'K':
			b.state.wOO = true
		case 'Q':
			b.state.wOOO = true
		case 'k':
			b.state.bOO = true
		case 'q':
			b.state.bOOO = true
		case '-':
			continue
		default:
			return nil, errors.New(fmt.Sprintf("bad castling char: %c", c))
		}
	}

	// Parse en passant target.
	if target, err := CoordFromString(parts[3]); err != nil {
		return nil, fmt.Errorf("error parsing en passant target: %w", err)
	} else {
		b.state.epTarget = target
	}

	// Parse the half move.
	if m, err := strconv.Atoi(parts[4]); err != nil {
		return nil, fmt.Errorf("error parsing half moves: %w", err)
	} else if m < 0 {
		return nil, fmt.Errorf("halfmove < 0: %d", m)
	} else {
		b.state.halfMove = m
	}

	// Parse the full move.
	if m, err := strconv.Atoi(parts[5]); err != nil {
		return nil, fmt.Errorf("error parsing full moves: %w", err)
	} else if m < 1 {
		return nil, fmt.Errorf("fullmove < 1: %d", m)
	} else {
		b.state.fullMove = m
	}

	if err := b.validate(); err != nil {
		return nil, err
	}

	// Figure out if the kings are in check.
	b.updateChecks()

	return b, nil
}

// CurrentPlayerScore returns the score for the current player.
func (b *Board) CurrentPlayerScore() Score {
	if b.state.turn == White {
		return b.state.score
	}
	return -b.state.score
}

// reversePlayers reverses the white and black players.
func (b *Board) reversePlayers() *Board {
	// Flip white/black
	hasSpace := false
	f := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			hasSpace = true
		}
		if hasSpace {
			switch r {
			case 'b':
				return 'w'
			case 'w':
				return 'b'
			}
		} else {
			switch {
			case unicode.IsUpper(r):
				return unicode.ToLower(r)
			case unicode.IsLower(r):
				return unicode.ToUpper(r)
			}
		}
		return r
	}, b.FENString())

	// Now we need to flip the board.
	pieces := strings.Split(f, " ")
	ranks := strings.Split(pieces[0], "/")
	for i := 0; i < len(ranks)/2; i++ {
		ranks[i], ranks[7-i] = ranks[7-i], ranks[i]
	}
	pieces[0] = strings.Join(ranks, "/")
	f = strings.Join(pieces, " ")

	// Finally, return the board.
	b2, err := FromFEN(f)
	if err != nil {
		panic(err)
	}
	return b2
}

func init() {
	r := rand.New(rand.NewSource(99))
	for i := range zLookups {
		zLookups[i] = Hash(r.Uint64())
	}
}
