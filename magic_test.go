package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestRookMagic(t *testing.T) {
	coord := testingCoordFunc(t)
	coords := func(str string) (coords []Coord) {
		for _, c := range strings.Split(str, " ") {
			coords = append(coords, coord(c))
		}
		return coords
	}
	tests := []struct {
		loc   Coord
		occ   Bit
		moves []Coord
	}{
		{coord("a1"), Bit(0x0000000000000001), coords("b1 c1 d1 e1 f1 g1 h1 a2 a3 a4 a5 a6 a7 a8")},
		{coord("a1"), Bit(0x0000000000000101), coords("b1 c1 d1 e1 f1 g1 h1 a2")},
	}

	for _, test := range tests {
		t.Logf("RookOccupancy(%v, %016x) %v", test.loc, test.occ.Uint64(), test.moves)
		locs := RookLookup(test.loc, test.occ)
		if !reflect.DeepEqual(locs, test.moves) {
			t.Errorf("RookOccupancy(%v, %016x) = %v, expected %v", test.loc, test.occ.Uint64(), locs, test.moves)
		}
	}
}