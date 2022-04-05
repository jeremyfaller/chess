package main

import "testing"

func TestIsCastle(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc     string
		p        Piece
		from, to Coord
		is       bool
	}{
		{"e4", White | Pawn, coord("e2"), coord("e4"), false},
		{"Ke2", White | King, coord("e1"), coord("e2"), false},
		{"O-O", White | King, coord("e1"), coord("g1"), true},
		{"O-O-O", White | King, coord("e1"), coord("c1"), true},
	}

	for _, test := range tests {
		move := Move{p: test.p, from: test.from, to: test.to}
		if v := move.IsCastle(); v != test.is {
			t.Errorf("[%v] IsCastle = %t, expected %t", test.desc, v, test.is)
		}

		if move.IsCastle() {
			if s := move.String(); s != test.desc {
				t.Errorf("[%v] String() = %q, expected %q", test.desc, s, test.desc)
			}
		}
	}
}

func TestIsPromotion(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc     string
		p        Piece
		from, to Coord
		is       bool
	}{
		{"Ke8", White | King, coord("e7"), coord("e8"), false},
		{"e8", White | Pawn, coord("e7"), coord("e8"), true},
		{"e1", Black | Pawn, coord("e2"), coord("e1"), true},
	}

	for _, test := range tests {
		move := Move{p: test.p, to: test.to, from: test.from}
		if v := move.IsPromotion(); v != test.is {
			t.Errorf("[%v] IsPromotion() = %t, expected %t", test.desc, v, test.is)
		}
	}
}
