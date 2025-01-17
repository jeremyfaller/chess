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
			if h := p.HashIdx(); h < 0 || h >= len(zLookups) {
				t.Errorf("[%v].HashIdx() = %d %d", p, h, p.Colorless()-Pawn)
			}
		}
	}
}

func TestPieceScore(t *testing.T) {
	// Test that all reasonable piece values have a score.
	for _, p := range []Piece{Empty, Pawn, Rook, Bishop, Knight, Queen, King} {
		if p != Empty && p.Score() == 0 {
			t.Errorf("%s.Score() == %d, expected != 0", p, p.Score())
		}
		// White should be positive.
		if p.Score() != (p | White).Score() {
			t.Errorf("%s.Score() = %d, expected %d", (p | White), (p | White).Score(), p.Score())
		}
		// And black is negative.
		if -p.Score() != (p | Black).Score() {
			t.Errorf("%s.Score() = %d, expected %d", (p | Black), (p | Black).Score(), -p.Score())
		}
	}
}

func TestPieceString(t *testing.T) {
	if str := scores[Pawn].String(); str != "1.00" {
		t.Errorf("%s != %s", str, "1.00")
	}
}

func TestPieceAttacks(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		p     Piece
		coord Coord
		occ   Bit
		sqs   []string
	}{
		{Piece(White | King), coord("a1"), 0, []string{"a2", "b2", "b1"}},
	}
	for i, test := range tests {
		var v Bit
		for _, s := range test.sqs {
			v.Set(coord(s).Idx())
		}
		if a := test.p.Attacks(test.coord, test.occ); a != v {
			t.Errorf("[%d] expected %v.Attacks(%v, %x) = %x, expected %x", i, test.p, test.coord, test.occ, a, v)
		}
	}
}
