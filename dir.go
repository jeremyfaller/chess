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

func (d Dir) String() string {
	switch d {
	case N:
		return "N"
	case NE:
		return "NE"
	case E:
		return "E"
	case SE:
		return "SE"
	case S:
		return "S"
	case SW:
		return "SW"
	case W:
		return "W"
	case NW:
		return "NW"
	case NNE:
		return "NNE"
	case NEE:
		return "NEE"
	case SEE:
		return "SEE"
	case SSE:
		return "SSE"
	case SSW:
		return "SSW"
	case SWW:
		return "SWW"
	case NWW:
		return "NWW"
	case NNW:
		return "NNW"
	case E2:
		return "E2"
	case W2:
		return "W2"
	case NN:
		return "NN"
	case SS:
		return "SS"
	}
	panic("unhandled direction")
}
