package main

import (
	"errors"
	"fmt"
	"math/bits"
	"testing"
	"unicode"
)

var noteString = "abcdefgh"
var UnknownDir = errors.New("bad dir")

type Bit uint64

func coordsFromBit(b Bit) []Coord {
	ret := make([]Coord, 0, bits.OnesCount64(uint64(b)))
	for i := 0; i < 64; i++ {
		sq := Bit(1) << i
		if b&sq != 0 {
			ret = append(ret, CoordFromBit(sq))
		}
	}
	return ret
}

type Coord int

var InvalidCoord = Coord(-1)

func CoordFromIdx(p int) Coord {
	if p < 0 || p > 63 {
		return InvalidCoord
	}
	return Coord(p)
}

func CoordFromBit(b Bit) Coord {
	return CoordFromIdx(bits.Len64(uint64(b)) - 1)
}

func (c Coord) IsValid() bool {
	return c >= 0 && c <= 63
}

func (c Coord) String() string {
	if !c.IsValid() {
		return "-"
	}
	return fmt.Sprintf("%c%d", noteString[c.X()], c.Y()+1)
}

func (c Coord) Idx() int {
	return int(c)
}

// Rank returns the rank of a Coord [0..7].
func (c Coord) Rank() int {
	return int(c / 8)
}

// Y returns the y cartesian value of a Coord.
func (c Coord) Y() int {
	return int(c / 8)
}

// File returns the file of a Coord[0..7].
func (c Coord) File() int {
	return int(c % 8)
}

// X returns the x cartesian value of a Coord.
func (c Coord) X() int {
	return int(c % 8)
}

// XDist returns the x-distance between two Coords.
func (c Coord) XDist(c2 Coord) int {
	return c.X() - c2.X()
}

// YDist returns the y-distance between two Coords.
func (c Coord) YDist(c2 Coord) int {
	return c.Y() - c2.Y()
}

// Bit returns a unique bit for a given Coord.
func (c Coord) Bit() Bit {
	return Bit(1) << c.Idx()
}

// CoordFromXY returns a Coord from an x/y pair.
func CoordFromXY(x, y int) Coord {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return InvalidCoord
	}
	return Coord(x + y*8)
}

func (c Coord) ApplyDir(d Dir) Coord {
	switch d {
	case N:
		return CoordFromXY(c.X(), c.Y()+1)
	case NE:
		return CoordFromXY(c.X()+1, c.Y()+1)
	case E:
		return CoordFromXY(c.X()+1, c.Y())
	case SE:
		return CoordFromXY(c.X()+1, c.Y()-1)
	case S:
		return CoordFromXY(c.X(), c.Y()-1)
	case SW:
		return CoordFromXY(c.X()-1, c.Y()-1)
	case W:
		return CoordFromXY(c.X()-1, c.Y())
	case NW:
		return CoordFromXY(c.X()-1, c.Y()+1)

	// knight moves
	case NNE:
		return CoordFromXY(c.X()+1, c.Y()+2)
	case NEE:
		return CoordFromXY(c.X()+2, c.Y()+1)
	case SEE:
		return CoordFromXY(c.X()+2, c.Y()-1)
	case SSE:
		return CoordFromXY(c.X()+1, c.Y()-2)
	case SSW:
		return CoordFromXY(c.X()-1, c.Y()-2)
	case SWW:
		return CoordFromXY(c.X()-2, c.Y()-1)
	case NWW:
		return CoordFromXY(c.X()-2, c.Y()+1)
	case NNW:
		return CoordFromXY(c.X()-1, c.Y()+2)

	// castle moves
	case E2:
		if c.X() != 4 {
			return InvalidCoord
		}
		return CoordFromXY(c.X()+2, c.Y())
	case W2:
		if c.X() != 4 {
			return InvalidCoord
		}
		return CoordFromXY(c.X()-2, c.Y())

	// pawn moves
	case NN:
		return CoordFromXY(c.X(), c.Y()+2)
	case SS:
		return CoordFromXY(c.X(), c.Y()-2)
	}
	panic("not handled")
}

func CoordFromString(s string) (Coord, error) {
	if s == "-" {
		return InvalidCoord, nil
	}
	if len(s) != 2 {
		return InvalidCoord, fmt.Errorf("invalid coord: %q", s)
	}
	if s[0] < 'a' || s[0] > 'h' {
		return InvalidCoord, fmt.Errorf("invalid coord: %q", s)
	}
	if !unicode.IsNumber(rune(s[1])) {
		return InvalidCoord, fmt.Errorf("invalid coord: %q", s)
	}
	if c := int(s[1] - '0'); c < 1 || c > 8 {
		return InvalidCoord, fmt.Errorf("invalid coord: %q", s)
	} else {
		return CoordFromXY(int(s[0]-'a'), int(s[1]-'1')), nil
	}
}

func testingCoordFunc(t *testing.T) func(string) Coord {
	return func(s string) Coord {
		c, err := CoordFromString(s)
		if err != nil {
			t.Fatalf("invalid coord: %v", err)
		}
		return c
	}
}

// FileIdx returns the file as an index (eg, [0..7]).
func (c Coord) FileIdx() int {
	if c == InvalidCoord {
		panic("file idx on an invalid coord")
	}
	return c.X()
}

// dirs is an array of LUTs, one for each square. For each of the squares, it
// maps a given coordinate to the direciton it takes to get there. Note, that
// it only returns valid standard directions, not special directions for pawns,
// and castling.
var dirs [64][64]Dir

// Set up dirs.
func init() {
	for i := 0; i < 64; i++ {
		c := CoordFromIdx(i)
		for _, d := range []Dir{N, NE, E, SE, S, SW, W, NW} {
			last := c
			for {
				last = last.ApplyDir(d)
				if !last.IsValid() {
					break
				}
				dirs[i][last.Idx()] = d
			}
		}
		for _, d := range []Dir{NNE, NEE, SEE, SSE, SSW, SWW, NWW, NNW} {
			if last := c.ApplyDir(d); last.IsValid() {
				dirs[i][last.Idx()] = d
			}
		}
	}
}

// DirBetween returns a Dir, given two coords.
//
// Returns InvalidDir if there's not mapping.
func DirBetween(from, to Coord) Dir {
	return dirs[from.Idx()][to.Idx()]
}
