package main

import (
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"unicode"
)

//go:generate go run zobrist_gen.go

const (
	StartingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type Spaces [64]Piece

// BoardState contains the state of the board that's Undoable.
type BoardState struct {
	spaces             Spaces
	turn               Piece // Whose turn is it?
	fullMove, halfMove int
	epTarget           Coord
	hash               Hash
	score              Score

	// State of the kings.
	wOO, wOOO, bOO, bOOO bool  // Can the white/black king castle kingside/queenside?
	wkLoc, bkLoc         Coord // Where are the kings located?
	isWCheck, isBCheck   bool

	// Is this space occupied?
	occ Bit
}

type Board struct {
	state    BoardState
	moves    []Move
	oldState []BoardState
}

func (b *Board) at(c Coord) Piece {
	return b.state.spaces[c.Idx()]
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
		b.state.occ.Clear(idx)
	} else {
		b.state.occ.Set(idx)
	}
	b.state.spaces[idx] = p
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
			print("|", b.at(CoordFromXY(x, y)).String())
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
			c := CoordFromXY(x, y)
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

// IsKingInCheck returns true if the given piece's color's king is in check.
func (b *Board) IsKingInCheck(p Piece) bool {
	if p == Empty {
		panic("empty")
	}
	if p.Color() == White {
		return b.state.isWCheck
	}
	return b.state.isBCheck
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
		b.state.isWCheck = v
	} else {
		b.state.isBCheck = v
	}
}

// isSquareAttacked returns true if a given square is attacked by the given color.
func (b *Board) isSquareAttacked(c Coord, color Piece) bool {
	if c == InvalidCoord {
		return false
	}

	bit := Bit(1 << c.Idx())
	occ := b.state.occ
	for v := occ; v != 0; {
		i := bits.TrailingZeros64(uint64(v))
		p := b.at(CoordFromIdx(i))
		v.Clear(i)
		if p.Color() != color {
			continue
		}
		if p.Colorless() == Queen {
			if (Bishop|p.Color()).Attacks(CoordFromIdx(i), occ)&bit != 0 {
				return true
			}
			if (Rook|p.Color()).Attacks(CoordFromIdx(i), occ)&bit != 0 {
				return true
			}
		} else {
			if p.Attacks(CoordFromIdx(i), occ)&bit != 0 {
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
				if b.at(CoordFromXY(m.to.X()-1, m.to.Y())) != Empty {
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
			dist := m.to.Y() - m.from.Y()
			if dist == 2 || dist == -2 {
				// Pawns can't move 2 spaces if it's not from the start location.
				if m.p.Color() == White && m.from.Y() != 1 ||
					m.p.Color() == Black && m.from.Y() != 6 {
					return false
				}

				// Pawns can't move through a piece.
				mid := CoordFromXY(m.from.X(), (m.from.Y()+m.to.Y())/2)
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
			m.isCapture = true
		} else if p2.Color() == m.p.Color() {
			// Cannot capture your own color.
			return false
		} else {
			m.isCapture = true
		}
	} else if p2 != Empty {
		// Can't take pieces of your own color.
		if p2.Color() == m.p.Color() {
			return false
		}
		// Can't take a king.
		if p2.Colorless() == King {
			return false
		}
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

	// A queen move is the same as a rook and bishop move.
	// Rather than a more complicated code structure involving
	// possible allocations, the cleanest way to do these two
	// piece check is with a goto.
	pCheck := p
	if p.Colorless() == Queen {
		pCheck = Bishop | p.Color()
	}

queenCheckRook:
	for _, toPos := range pCheck.Moves(c, b.state.occ) {
		// If we'd overlap one our own pieces, skip it.
		if b.at(toPos).Color() == p.Color() {
			continue
		}

		move := Move{
			p:    p,
			to:   toPos,
			from: c,
		}

		// If the move would be illegal, we keep checking as a different move in
		// this direction might be legal.
		if !b.isLegalMove(&move) {
			continue
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
	}

	// See above for why we use a goto here.
	if p.Colorless() == Queen && pCheck.Colorless() == Bishop {
		pCheck = Rook | p.Color()
		goto queenCheckRook
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
	distY := (m.to.Y() - m.from.Y())
	if m.p.IsPawn() && (distY == 2 || distY == -2) {
		return CoordFromXY(m.to.X(), (m.to.Y()+m.from.Y())/2)
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

	if p.IsWhite() && c.Y() == 0 {
		if c.X() == 0 {
			b.state.wOOO = false
		} else if c.X() == 7 {
			b.state.wOO = false
		}
	} else if p.IsBlack() && c.Y() == 7 {
		if c.X() == 0 {
			b.state.bOOO = false
		} else if c.X() == 7 {
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
			dist := m.to.X() - m.from.X()
			rookFrom, rookTo := CoordFromXY(7, m.from.Y()), CoordFromXY(5, m.from.Y())
			if dist < 0 { // queenside
				rookFrom = CoordFromXY(0, rookFrom.Y())
				rookTo = CoordFromXY(3, rookTo.Y())
			}
			b.set(Empty, rookFrom)
			b.set(m.p.Color()|Rook, rookTo)
		}
		if m.isEnPassant { // Need to remove captured pawn.
			c := m.to
			if m.p.Color() == White {
				c = CoordFromXY(c.X(), 4)
			} else {
				c = CoordFromXY(c.X(), 3)
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
	b.moves = b.moves[:len(b.moves)-1]
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
	counts := make(map[perftHash]uint64, 1000000)

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
				counts[h] = cnt
			}

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
		oldState: make([]BoardState, 0, 200),
		moves:    make([]Move, 0, 200),
	}
}

// New returns a new Board, set up for play (ie a new chess game).
func New() *Board {
	b, err := FromFEN(StartingFEN)
	if err != nil {
		panic(err)
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
	coord := CoordFromXY(0, 7)
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
		coord = CoordFromXY(
			(coord.X()+inc)%8,
			(coord.Y() - (coord.X()+inc)/8),
		)
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
