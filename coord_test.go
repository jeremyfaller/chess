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
