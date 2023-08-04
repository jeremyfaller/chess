package main

import "fmt"

type Move struct {
	p           Piece
	from, to    Coord
	promotion   Piece
	isCapture   bool
	isEnPassant bool // true if an en passant capture.
	isCheck     bool // true if move is check.
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
		return m.to.Y() == 7
	}
	return m.to.Y() == 0
}

// IsCastle returns true if the move would be a castling move.
func (m *Move) IsCastle() bool {
	if !m.p.IsKing() {
		return false
	}
	xDist := m.to.XDist(m.from)
	return xDist == 2 || xDist == -2
}

// IsKingsideCastle returns true if a Move is a king side castle.
func (m *Move) IsKingsideCastle() bool {
	if !m.IsCastle() {
		return false
	}
	return m.to.X() > m.from.X()
}

// IsQueensideCastle returns true if a Move is a queen side castle.
func (m *Move) IsQueensideCastle() bool {
	if !m.IsCastle() {
		return false
	}
	return m.to.X() < m.from.X()
}

// CastleMidCoord returns the Coord for the middle of a castling move.
func (m *Move) CastleMidCoord() Coord {
	if !m.IsCastle() {
		return InvalidCoord
	}
	return CoordFromXY((m.to.X()+m.from.X())/2, m.to.Y())
}

// RookCoord returns the Coord of the when it's a castling move.
func (m *Move) RookCoord() Coord {
	if m.p.IsKing() {
		xDist := m.to.XDist(m.from)
		if xDist == 2 {
			return CoordFromXY(7, m.to.Y())
		} else if xDist == -2 {
			return CoordFromXY(0, m.to.Y())
		}
	}
	return InvalidCoord
}

func (m Move) castleString() string {
	if !m.IsCastle() {
		panic("not a castling move")
	}
	xDist := m.to.XDist(m.from)
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

// promotionString returns the string of the promotion.
func (m Move) promotionString() string {
	if !m.IsPromotion() {
		return ""
	}
	return "=" + m.promotion.NoteString()
}

// checkString returns the string if we're in check.
func (m Move) checkString() string {
	if m.isCheck {
		return "+"
	}
	return ""
}

// algebraicString returns the move's string in algebraic notation.
func (m Move) algebraicString() string {
	if m.IsCastle() {
		return m.castleString()
	}
	var capString string
	if m.isCapture {
		capString = "x"
	}
	return fmt.Sprintf("%s%s%s%s%s", m.from.String(), capString, m.to.String(), m.promotionString(), m.checkString())
}

// figureString returns the move's string in figure notation.
func (m Move) figureString() string {
	var capString string
	if m.isCapture {
		capString = "x"
	}
	return fmt.Sprintf("%s%s%s%s%s", m.p.NoteString(), capString, m.to.String(), m.promotionString(), m.checkString())
}

// longAlgrbraicString returns long algebraic moves.
func (m Move) longAlgebraicString() string {
	promo := ""
	if m.IsPromotion() {
		promo = m.promotion.NoteString()
	}
	return fmt.Sprintf("%s%s%s", m.from.String(), m.to.String(), promo)
}

// String returns a string for the given Move. Note that it doesn't handle
// ambiguous moves, eg Nef4.
func (m Move) String() string {
	// Special-case castling.
	if m.IsCastle() {
		return m.castleString()
	}
	return m.longAlgebraicString()
}

// IsVertical returns true if a Move is only a vertical move.
func (m Move) IsVertical() bool {
	return m.to.X() == m.from.X()
}

// IsHorizontal returns true if a Move is a horizontal move.
func (m Move) IsHorizontal() bool {
	return m.to.Y() == m.from.Y()
}
