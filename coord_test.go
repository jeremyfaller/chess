package main

import "testing"

func TestCoordConversion(t *testing.T) {
	bits := make(map[uint64]struct{})
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			e := Coord{x, y}
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

func TestDirForCoord(t *testing.T) {
	for p := 0; p < 64; p++ {
		c := CoordFromIdx(p)
		for _, d := range []Dir{N, NE, E, SE, S, SW, W, NW} {
			last := c
			for {
				last = last.ApplyDir(d)
				if !last.IsValid() {
					break
				}
				if dir := DirBetween(c, last); dir != d {
					t.Errorf("[%v, dir:%v] DirBetween(%v, %v) = %v, expected = %v", c, d, c, last, dir, d)
				}
			}
		}
		for _, d := range []Dir{NNE, NEE, SEE, SSE, SSW, SWW, NWW, NNW} {
			if last := c.ApplyDir(d); last.IsValid() {
				if dir := DirBetween(c, last); dir != d {
					t.Errorf("[%v, dir:%v] DirBetween(%v, %v) = %v, expected = %v", c, d, c, last, dir, d)
				}
			}
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
