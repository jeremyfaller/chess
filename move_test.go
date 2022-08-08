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

func TestIsKingsideCastle(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc string
		m    Move
		is   bool
	}{
		{"white pawn", Move{p: White | Pawn, from: coord("a2"), to: coord("a4")}, false},
		{"white king move", Move{p: White | King, from: coord("e1"), to: coord("e2")}, false},
		{"white king castle", Move{p: White | King, from: coord("e1"), to: coord("g1")}, true},
		{"black king castle", Move{p: Black | King, from: coord("e1"), to: coord("g1")}, true},
		{"white queen castle", Move{p: White | King, from: coord("e1"), to: coord("c1")}, false},
		{"black queen castle", Move{p: Black | King, from: coord("e1"), to: coord("c1")}, false},
	}

	for i, test := range tests {
		if test.m.IsKingsideCastle() != test.is {
			t.Errorf("[%d] %v IsKingsideCastle = %v, expected %v", i, test.desc, !test.is, test.is)
		}
	}
}

func TestIsQueensideCastle(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc string
		m    Move
		is   bool
	}{
		{"white pawn", Move{p: White | Pawn, from: coord("a2"), to: coord("a4")}, false},
		{"white king move", Move{p: White | King, from: coord("e1"), to: coord("e2")}, false},
		{"white king castle", Move{p: White | King, from: coord("e1"), to: coord("g1")}, false},
		{"black king castle", Move{p: Black | King, from: coord("e1"), to: coord("g1")}, false},
		{"white queen castle", Move{p: White | King, from: coord("e1"), to: coord("c1")}, true},
		{"black queen castle", Move{p: Black | King, from: coord("e1"), to: coord("c1")}, true},
	}

	for i, test := range tests {
		if test.m.IsQueensideCastle() != test.is {
			t.Errorf("[%d] %v IsQueensideCastle = %v, expected %v", i, test.desc, !test.is, test.is)
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

func TestPromotionString(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc string
		move Move
		ex   string
	}{
		{"white queen", Move{p: White | Pawn, from: coord("e7"), to: coord("e8"), promotion: White | Queen}, "e7e8=Q"},
		{"white queen", Move{p: White | Pawn, from: coord("e7"), to: coord("d8"), promotion: White | Queen, isCapture: true}, "e7xd8=Q"},
		{"black knight", Move{p: Black | Pawn, from: coord("e2"), to: coord("e1"), promotion: Black | Knight}, "e2e1=n"},
	}
	for i, test := range tests {
		if str := test.move.String(); str != test.ex {
			t.Errorf("[%d] %s %v = %q, expected %q", i, test.desc, test.move, str, test.ex)
		}
	}
}
