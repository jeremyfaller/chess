package main

import (
	"fmt"
	"testing"
)

func TestCoordConversion(t *testing.T) {
	bits := make(map[uint64]struct{})
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			e := CoordFromXY(x, y)
			ep := e.Idx()
			eb := e.Bit()
			if c := CoordFromIdx(ep); c != e {
				t.Errorf("[%v] CoordFromPos(%v) = %v, expected %v", e, ep, c, e)
			}
			if c := CoordFromBit(eb); c != e {
				t.Errorf("[%v] CoordFromBit(%v) = %v, expected %v", e, ep, c, e)
			}
			bits[uint64(eb)] = struct{}{}
		}
	}
	if len(bits) != 64 { // all bits should be unique
		t.Errorf("coords not unique")
	}
}

func TestInvalidCoord(t *testing.T) {
	o := CoordFromXY(0, 0)
	l := CoordFromIdx(63)
	tests := []Coord{
		CoordFromXY(o.X()-1, o.Y()),
		CoordFromXY(o.X(), o.Y()-1),
		CoordFromXY(o.X()-1, o.Y()-1),
		CoordFromXY(l.X()+1, l.Y()),
		CoordFromXY(l.X(), l.Y()+1),
		CoordFromXY(l.X()+1, l.Y()+1),
		CoordFromXY(-1, -1),
	}
	for i, test := range tests {
		if test.IsValid() {
			t.Errorf("[%d] %v.IsValid() = true", i, test)
		}
	}
}

func TestCoordIndices(t *testing.T) {
	x, y := "a", 1
	for i := 0; i < 64; i++ {
		c := CoordFromIdx(i)
		s := fmt.Sprintf("%s%d", x, y)
		if c.String() != s {
			t.Errorf("idx: %d – %s != %s", i, c.String(), s)
		}
		x = fmt.Sprintf("%c", x[0]+1)
		if x == "i" {
			x = "a"
			y += 1
		}
	}
}

func TestBadCastleStartingLocations(t *testing.T) {
	tests := []struct {
		desc  string
		dir   Dir
		loc   Coord
		valid bool
	}{
		{"a1,O-O", E2, CoordFromIdx(0), false},
		{"a1,O-O-O", W2, CoordFromIdx(0), false},
		{"e1,O-O", E2, CoordFromIdx(4), true},
		{"e1,O-O", W2, CoordFromIdx(4), true},
		{"e8,O-O", E2, CoordFromIdx(8*7 + 4), true},
		{"e8,O-O-O", W2, CoordFromIdx(8*7 + 4), true},
		{"c8,O-O-O", W2, CoordFromIdx(8*7 + 2), false},
	}
	for i, test := range tests {
		t.Logf("%s\n", test.desc)
		v := test.loc.ApplyDir(test.dir)
		if value := v.IsValid(); value != test.valid {
			t.Errorf("[%d] %s IsValid() = %v, expected %v", i, test.desc, value, test.valid)
		}
	}
}
