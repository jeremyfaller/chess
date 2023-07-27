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
		t.Logf("RookOccupancy(%v, %v) %v", test.loc, test.occ, test.moves)
		locs := rookLookup(test.loc, test.occ)
		if !reflect.DeepEqual(locs, test.moves) {
			t.Errorf("RookOccupancy(%v, %v) = %v, expected %v", test.loc, test.occ, locs, test.moves)
		}
	}
}

func TestMagicEquivalence(t *testing.T) {
	cmpBits := func(name string, coords [64][][]Coord, bits [64][]Bit) {
		for i := range bits {
			if len(bits[i]) != len(coords[i]) {
				t.Fatalf("len(%sBits[%d]) != len(%sCoords[%d])", name, i, name, i)
			}
			for j := range coords[i] {
				if val := ToBit(coords[i][j]); bits[i][j] != val {
					t.Fatalf("%sBits[%d][%d] = %v, expected %v – %v %v", name, i, j, bits[i][j], val,
						coords[i][j], val)
				}
			}
		}
	}
	cmpBits("rook", rookCoords, rookBits)
	cmpBits("bishop", bishopCoords, bishopBits)
}
