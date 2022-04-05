package main

import "testing"

func TestOppositeColor(t *testing.T) {
	tests := []struct {
		p, op Piece
	}{
		{White | King, Black},
		{Black | King, White},
	}

	for _, test := range tests {
		if c := test.p.OppositeColor(); c != test.op {
			t.Errorf("%v.OppositeColor() = %v, expected = %v", test.p, c, test.op)
		}
	}
}

func TestHashIdx(t *testing.T) {
	for _, c := range []Piece{White, Black} {
		for _, ty := range []Piece{Pawn, Rook, Knight, Bishop, Queen, King} {
			p := c | ty
			if h := p.HashIdx(); h < 0 || h >= 12 {
				t.Errorf("[%v].HashIdx() = %d %d", p, h, p.Colorless()-Pawn)
			}
		}
	}
}
