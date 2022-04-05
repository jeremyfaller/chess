package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
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
func (a *PsuedoMoves) Attackers(c Coord) []Coord {
	return a[c.Idx()].Attackers()
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

type Board struct {
	spaces    Spaces
	hash      Hash
	epTarget  Coord
	turn      Piece
	wOO, wOOO bool // can white castle kingside/queenside?
	bOO, bOOO bool // can black castle kingside/queenside?
	halfMove  int
	fullMove  int
	moves     []Move

	// State of the kings.
	wkLoc, bkLoc       Coord
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

	// Update the hash.
	hashP := p
	if hashP == Empty {
		hashP = b.at(c)
	}
	if hashP != Empty {
		b.hash ^= zLookups[hashP.HashIdx()*idx]
	}

	// Update the king's location.
	if p.Colorless() == King {
		if p.Color() == White {
			b.wkLoc = c
		} else {
			b.bkLoc = c
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
	return b.hash
}

// Print pretty-prints a board on standard out.
func (b *Board) Print() {
	for y := 7; y >= 0; y-- {
		println("+-+-+-+-+-+-+-+-+")
		for x := 0; x < 8; x++ {
			print("|", b.at(Coord{x, y}).String())
		}
		println("|")
	}
	println("+-+-+-+-+-+-+-+-+")
}

func (b *Board) castleString() string {
	var s string
	canCastle := false
	if b.wOO {
		canCastle = true
		s += "K"
	}
	if b.wOOO {
		canCastle = true
		s += "Q"
	}
	if b.bOO {
		canCastle = true
		s += "k"
	}
	if b.bOOO {
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

	// Print the board
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
	if b.turn == White {
		s += "w "
	} else {
		s += "b "
	}

	// Print the Casting rights
	s += b.castleString() + " "

	// Print the rest.
	s += fmt.Sprintf("%v %d %d", b.epTarget, b.halfMove, b.fullMove)
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
		return b.wkLoc
	}
	return b.bkLoc
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
func (b *Board) doesSquareAttack(from, to Coord) bool {
	p := b.at(from)
	// If there's no piece at the attacking square, clearly no.
	if p == Empty {
		return false
	}

	// Check that the colors are different.
	if b.at(to).Color() == p.Color() {
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

	for _, at := range b.PseudoMoves(color).Attackers(c) {
		if b.doesSquareAttack(at, c) {
			return true
		}
	}
	return false
}

// wouldKingBeInCheck returns true if a move would result in an illegal check.
func (b *Board) wouldKingBeInCheck(m *Move) bool {
	// If the king's not in check, and none of the pieces that attack the from
	// square also attack the king, the move will not result in a check.
	isInCheck := b.IsKingInCheck(m.p)
	kloc := b.KingLoc(m.p)
	attacked := b.PseudoMoves(m.p.OppositeColor())
	bits := attacked[kloc.Idx()] & attacked[m.from.Idx()]
	if !isInCheck && bits == 0 {
		return false
	}

	// So, if the king was in check, we need to see if we would block the
	// check, or take the checking piece.
	b.MakeMove(*m)
	check := b.IsKingInCheck(m.p)
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
			// Can't castle out of check.
			if b.IsKingInCheck(m.p) {
				return false
			}

			mid := m.CastleMidCoord()

			// Can't castle across or into non-empty squares.
			if p2 != Empty || b.at(mid) != Empty {
				return false
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
			if m.to != b.epTarget {
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

func (b *Board) addMoves(c Coord, moves []Move) []Move {
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

			// Check to see if the move is leagl.
			move := Move{
				p:    p,
				to:   toPos,
				from: c,
			}

			// If the move isn't legal, none of the rest in this direction will be.
			if !b.isLegalMove(&move) {
				break
			}

			// It's a legal move.
			lastPos = toPos
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
func (b *Board) PossibleMoves() []Move {
	moves := []Move{}
	for p := 0; p < 64; p++ {
		c := CoordFromIdx(p)
		p := b.at(c)
		if p.IsEmpty() { // no piece, skip
			continue
		}
		if p.Color() != b.turn { // not the color we want, skip
			continue
		}
		moves = b.addMoves(c, moves)
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
			b.wOOO = false
		} else if c.x == 7 {
			b.wOO = false
		}
	} else if p.IsBlack() && c.y == 7 {
		if c.x == 0 {
			b.bOOO = false
		} else if c.x == 7 {
			b.bOO = false
		}
	}
}

// zCastle returns the Zobrist hash for the given castle state.
func (b *Board) zCastle() (v Hash) {
	if b.wOO {
		v ^= zWOO
	}
	if b.wOOO {
		v ^= zWOOO
	}
	if b.bOO {
		v ^= zBOO
	}
	if b.bOOO {
		v ^= zBOOO
	}
	return v
}

// updateCastleState updates the castling state for a given move.
func (b *Board) updateCastleState(m Move) {
	b.hash ^= b.zCastle()
	// If we're moving a king, we can't castle anymore.
	if m.p.IsKing() {
		if m.p.IsWhite() {
			b.wOO, b.wOOO = false, false
		} else {
			b.bOO, b.bOOO = false, false
		}
	}

	// If we're moving or capturing a rook, update the state.
	b.handleRookMoveOrCapture(m.from)
	b.handleRookMoveOrCapture(m.to)
	b.hash ^= b.zCastle()
}

func (b *Board) updateEPTarget(m Move) {
	if b.epTarget != InvalidCoord {
		b.hash ^= zLookups[zEP+b.epTarget.FileIdx()]
	}
	b.epTarget = epTarget(m)
	if b.epTarget != InvalidCoord {
		b.hash ^= zLookups[zEP+b.epTarget.FileIdx()]
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
	m.prevHash = b.hash
	m.prevEP, m.prevFull, m.prevHalf = b.epTarget, b.fullMove, b.halfMove
	m.prevWOO, m.prevWOOO, m.prevBOO, m.prevBOOO = b.wOO, b.wOOO, b.bOO, b.bOOO

	// Update turn variables, and board state.
	b.hash ^= zLookups[zBlack]
	if b.turn == White {
		b.turn = Black
	} else {
		b.fullMove += 1
		b.turn = White
	}
	b.halfMove += 1
	if m.p.IsPawn() || m.isCapture {
		b.halfMove = 0
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

	// Fix game state.
	b.turn = m.p.Color()
	b.fullMove, b.halfMove = m.prevFull, m.prevHalf
	b.epTarget = m.prevEP
	b.wOO, b.wOOO, b.bOO, b.bOOO = m.prevWOO, m.prevWOOO, m.prevBOO, m.prevBOOO

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
	b.hash = m.prevHash

	b.updateChecks()
}

// GetMove gets a move given two coordinates.
func (b *Board) GetMove(from, to Coord) (Move, error) {
	if !from.IsValid() {
		return Move{}, errors.New("invalid from coord")
	}
	if !to.IsValid() {
		return Move{}, errors.New("invalid to coord")
	}
	for _, m := range b.addMoves(from, nil) {
		if m.from != from || m.to != to {
			continue
		}
		return m, nil
	}
	return Move{}, fmt.Errorf("invalid move: %v", Move{from: from, to: to})
}

// Perft calcualtes the number of possible moves at a given depth. It's quite
// helpful debugging the move generation. Optionally, Perft will also print the
// number of reachable moves for each valid move in the given board state.
func (b Board) Perft(depth int, showMoves bool) int {
	var perft func(int, bool) int

	perft = func(d int, s bool) int {
		if d == 0 {
			return 1
		}
		total := 0
		for _, move := range b.PossibleMoves() {
			b.MakeMove(move)
			cnt := perft(d-1, false)
			if s {
				fmt.Printf("%v: %d\n", move, cnt)
			}
			b.UnmakeMove()
			total += cnt
		}
		return total
	}
	return perft(depth, showMoves)
}

// EmptyBoard returns a new, empty board. No state of gameplay is set up.
func EmptyBoard() *Board {
	return &Board{
		epTarget: InvalidCoord,
		fullMove: 1,
		wkLoc:    InvalidCoord,
		bkLoc:    InvalidCoord,
	}
}

// New returns a new Board, set up for play (ie a new chess game).
func New() *Board {
	b := EmptyBoard()
	b.turn = White
	b.wOO = true
	b.wOOO = true
	b.bOO = true
	b.bOOO = true
	pieces := []Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for x, p := range pieces {
		b.set(p|White, Coord{x, 0})
		b.set(Pawn|White, Coord{x, 1})
		b.set(Pawn|Black, Coord{x, 6})
		b.set(p|Black, Coord{x, 7})
	}
	return b
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
	if b.wkLoc == InvalidCoord {
		return errors.New("no white king on board")
	}
	if b.bkLoc == InvalidCoord {
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
	b.turn = White
	if parts[1] == "b" {
		b.turn = Black
	}

	// Parse Castling.
	for _, c := range parts[2] {
		switch c {
		case 'K':
			b.wOO = true
		case 'Q':
			b.wOOO = true
		case 'k':
			b.bOO = true
		case 'q':
			b.bOOO = true
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
		b.epTarget = target
	}

	// Parse the half move.
	if m, err := strconv.Atoi(parts[4]); err != nil {
		return nil, fmt.Errorf("error parsing half moves: %w", err)
	} else if m < 0 {
		return nil, fmt.Errorf("halfmove < 0: %d", m)
	} else {
		b.halfMove = m
	}

	// Parse the full move.
	if m, err := strconv.Atoi(parts[5]); err != nil {
		return nil, fmt.Errorf("error parsing full moves: %w", err)
	} else if m < 1 {
		return nil, fmt.Errorf("fullmove < 1: %d", m)
	} else {
		b.fullMove = m
	}

	if err := b.validate(); err != nil {
		return nil, err
	}

	// Figure out if the kings are in check.
	b.updateChecks()

	return b, nil
}

func init() {
	r := rand.New(rand.NewSource(99))
	for i := range zLookups {
		zLookups[i] = Hash(r.Uint64())
	}
}
