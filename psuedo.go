package main

import "fmt"

type PsuedoMoves [64]Bit

// Set sets a position as attacked.
func (a *PsuedoMoves) Update(p Piece, c Coord) {
	// If we're clearing a piece, just clear all spaces on the board it attacks.
	if p == Empty {
		bit := ^c.Bit()
		for i := 0; i < 64; i += 8 {
			a[i+0] &= bit
			a[i+1] &= bit
			a[i+2] &= bit
			a[i+3] &= bit
			a[i+4] &= bit
			a[i+5] &= bit
			a[i+6] &= bit
			a[i+7] &= bit
		}
		return
	}

	// Set the attacked bits for the given piece.
	attacks := p.Psuedos(c)
	for i := 0; i < 64; i += 8 {
		a[i+0] |= (*attacks)[i+0]
		a[i+1] |= (*attacks)[i+1]
		a[i+2] |= (*attacks)[i+2]
		a[i+3] |= (*attacks)[i+3]
		a[i+4] |= (*attacks)[i+4]
		a[i+5] |= (*attacks)[i+5]
		a[i+6] |= (*attacks)[i+6]
		a[i+7] |= (*attacks)[i+7]
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
