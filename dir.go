package main

type Dir int

const (
	InvalidDir Dir = iota
	N
	NE
	E
	SE
	S
	SW
	W
	NW

	// Knight directions
	NNE
	NEE
	SEE
	SSE
	SSW
	SWW
	NWW
	NNW

	// Castling
	E2 // OO
	W2 // OOO

	// Pawn double moves
	NN
	SS
)

// IsKnight returns true if we're dealing with a knight move.
func (d Dir) IsKnight() bool {
	switch d {
	case NNE, NEE, SEE, SSE, SSW, SWW, NWW, NNW:
		return true
	default:
		return false
	}
}
