package main

import "fmt"

type Move struct {
	p           Piece
	from, to    Coord
	isCapture   bool
	isEnPassant bool // true if an en passant capture
	promotion   Piece

	// State for unwinding a move
	captured Piece
}

// Color returns the color of the move.
func (m *Move) Color() Piece {
	return m.p.Color()
}

// IsPromotion returns true if the move would be a promotion – it has nothing
// to do with the promotion field.
func (m *Move) IsPromotion() bool {
	if !m.p.IsPawn() {
		return false
	}
	if m.p.IsWhite() {
		return m.to.y == 7
	}
	return m.to.y == 0
}

// IsCastle returns true if the move would be a castling move.
func (m *Move) IsCastle() bool {
	if !m.p.IsKing() {
		return false
	}
	xDist := m.to.x - m.from.x
	return xDist == 2 || xDist == -2
}

// IsKingsideCastle returns true if a Move is a king side castle.
func (m *Move) IsKingsideCastle() bool {
	if !m.IsCastle() {
		return false
	}
	return m.to.x > m.from.x
}

// IsQueensideCastle returns true if a Move is a queen side castle.
func (m *Move) IsQueensideCastle() bool {
	if !m.IsCastle() {
		return false
	}
	return m.to.x < m.from.x
}

// CastleMidCoord returns the Coord for the middle of a castling move.
func (m *Move) CastleMidCoord() Coord {
	if !m.IsCastle() {
		return InvalidCoord
	}
	return Coord{(m.to.x + m.from.x) / 2, m.to.y}
}

// RookCoord returns the Coord of the when it's a castling move.
func (m *Move) RookCoord() Coord {
	if m.p.IsKing() {
		xDist := m.to.x - m.from.x
		if xDist == 2 {
			return Coord{x: 7, y: m.to.y}
		} else if xDist == -2 {
			return Coord{x: 0, y: m.to.y}
		}
	}
	return InvalidCoord
}

func (m Move) castleString() string {
	if !m.IsCastle() {
		panic("not a castling move")
	}
	xDist := m.to.x - m.from.x
	if m.p.IsWhite() {
		if xDist > 0 {
			return "O-O"
		}
		return "O-O-O"
	}
	if xDist > 0 {
		return "O-O-O"
	}
	return "O-O"
}

// algebraicString returns the move's string in algebraic notation.
func (m Move) algebraicString() string {
	if m.IsCastle() {
		return m.castleString()
	}
	var capString string
	if m.isCapture {
		//capString = "x"
	}
	return fmt.Sprintf("%s%s%s", m.from.String(), capString, m.to.String())
}

// figureString returns the move's string in figure notation.
func (m Move) figureString() string {
	var capString string
	if m.isCapture {
		capString = "x"
	}
	var promString string
	if m.IsPromotion() {
		promString = "=" + m.promotion.String()
	}
	return fmt.Sprintf("%s%s%s%s", m.p.NoteString(), capString, m.to.String(), promString)
}

// String returns a string for the given Move. Note that it doesn't handle
// ambiguous moves, eg Nef4.
func (m Move) String() string {
	// Special-case castling.
	if m.IsCastle() {
		return m.castleString()
	}
	return m.algebraicString()
}

// IsVertical returns true if a Move is only a vertical move.
func (m Move) IsVertical() bool {
	return m.to.x == m.from.x
}

// IsHorizontal returns true if a Move is a horizontal move.
func (m Move) IsHorizontal() bool {
	return m.to.y == m.from.y
}
