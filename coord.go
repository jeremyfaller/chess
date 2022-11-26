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

type Coord struct {
	x, y int
}

var InvalidCoord = Coord{-1, -1}

func CoordFromIdx(p int) Coord {
	return Coord{p % 8, p / 8}
}

func CoordFromBit(b Bit) Coord {
	return CoordFromIdx(bits.Len64(uint64(b)) - 1)
}

func (c Coord) IsValid() bool {
	return c.x >= 0 && c.x <= 7 && c.y >= 0 && c.y <= 7
}

func (c Coord) String() string {
	if !c.IsValid() {
		return "-"
	}
	return fmt.Sprintf("%c%d", noteString[c.x], c.y+1)
}

func (c Coord) PanicInvalid() {
	if !c.IsValid() {
		panic(fmt.Sprintf("bad coord: (%d,%d)", c.x, c.y))
	}
}

func (c Coord) Idx() int {
	c.PanicInvalid()
	return c.x + c.y*8
}

// Bit returns a unique bit for a given Coord.
func (c Coord) Bit() Bit {
	c.PanicInvalid()
	return Bit(1) << c.Idx()
}

func (c Coord) ApplyDir(d Dir) Coord {
	switch d {
	case N:
		return Coord{c.x, c.y + 1}
	case NE:
		return Coord{c.x + 1, c.y + 1}
	case E:
		return Coord{c.x + 1, c.y}
	case SE:
		return Coord{c.x + 1, c.y - 1}
	case S:
		return Coord{c.x, c.y - 1}
	case SW:
		return Coord{c.x - 1, c.y - 1}
	case W:
		return Coord{c.x - 1, c.y}
	case NW:
		return Coord{c.x - 1, c.y + 1}

	// knight moves
	case NNE:
		return Coord{c.x + 1, c.y + 2}
	case NEE:
		return Coord{c.x + 2, c.y + 1}
	case SEE:
		return Coord{c.x + 2, c.y - 1}
	case SSE:
		return Coord{c.x + 1, c.y - 2}
	case SSW:
		return Coord{c.x - 1, c.y - 2}
	case SWW:
		return Coord{c.x - 2, c.y - 1}
	case NWW:
		return Coord{c.x - 2, c.y + 1}
	case NNW:
		return Coord{c.x - 1, c.y + 2}

	// castle moves
	case E2:
		if c.x != 4 {
			return InvalidCoord
		}
		return Coord{c.x + 2, c.y}
	case W2:
		if c.x != 4 {
			return InvalidCoord
		}
		return Coord{c.x - 2, c.y}

	// pawn moves
	case NN:
		return Coord{c.x, c.y + 2}
	case SS:
		return Coord{c.x, c.y - 2}
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
		return Coord{int(s[0] - 'a'), int(s[1] - '1')}, nil
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
	return c.x
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
