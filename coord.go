package main

import (
	"errors"
	"fmt"
	"math/bits"
	"testing"
	"unicode"
)

//go:generate go run coord_gen.go coord.go dir.go bit.go

var noteString = "abcdefgh"
var UnknownDir = errors.New("bad dir")

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
	return int(c) & 0x3F
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
	return CoordFromIdx(x + y*8)
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

// ToBit takes a slice of Coord, and returns a Bit.
func ToBit(c []Coord) (b Bit) {
	for _, v := range c {
		b.Set(v.Idx())
	}
	return b
}

// TooCoordSlice returns a slice of Coord from a given Bit.
func (b Bit) ToCoordSlice() (c []Coord) {
	c = make([]Coord, 0, b.CountOnes())
	for b != 0 {
		c = append(c, b.NextCoord())
	}
	return c
}

func (b *Bit) NextCoord() Coord {
	i := 63 - bits.LeadingZeros64(uint64(*b))
	b.Clear(i)
	return CoordFromIdx(i)
}
