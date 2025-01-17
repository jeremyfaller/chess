// Generates magic bitboards.
//
// This code is adapated from the C code at:
//
// https://www.chessprogramming.org/index.php?title=Looking_for_Magics&oldid=2272

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"math/bits"
	"math/rand"
	"os"
	"sync"
)

func doesBlock(f, r int, block Bit) bool {
	return block&(Bit(1)<<(f+r*8)) != 0
}

func genRookMask(idx int) Bit {
	row := Bit(0x7E << (8 * (idx / 8)))
	col := Bit(0x0001010101010100 << (idx % 8))
	b := row | col
	b.Clear(idx)
	return b
}

func genRookAttacks(sq int, block Bit) (b Bit) {
	rk, fl := sq/8, sq%8
	for r := rk + 1; r <= 7; r++ {
		b |= (Bit(1) << (fl + r*8))
		if doesBlock(fl, r, block) {
			break
		}
	}
	for r := rk - 1; r >= 0; r-- {
		b |= (Bit(1) << (fl + r*8))
		if doesBlock(fl, r, block) {
			break
		}
	}
	for f := fl + 1; f <= 7; f++ {
		b |= (Bit(1) << (f + rk*8))
		if doesBlock(f, rk, block) {
			break
		}
	}
	for f := fl - 1; f >= 0; f-- {
		b |= (Bit(1) << (f + rk*8))
		if doesBlock(f, rk, block) {
			break
		}
	}
	return b
}

func genBishopMask(idx int) (b Bit) {
	for x, y := idx%8+1, idx/8+1; x < 7 && y < 7; x, y = x+1, y+1 {
		b.Set(x + y*8)
	}
	for x, y := idx%8-1, idx/8+1; x > 0 && y < 7; x, y = x-1, y+1 {
		b.Set(x + y*8)
	}
	for x, y := idx%8-1, idx/8-1; x > 0 && y > 0; x, y = x-1, y-1 {
		b.Set(x + y*8)
	}
	for x, y := idx%8+1, idx/8-1; x < 7 && y > 0; x, y = x+1, y-1 {
		b.Set(x + y*8)
	}
	return b
}

func genBishopAttacks(sq int, block Bit) (b Bit) {
	rk, fl := sq/8, sq%8
	for r, f := rk+1, fl+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		b |= Bit(1) << (f + r*8)
		if doesBlock(f, r, block) {
			break
		}
	}
	for r, f := rk+1, fl-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		b |= Bit(1) << (f + r*8)
		if doesBlock(f, r, block) {
			break
		}
	}
	for r, f := rk-1, fl+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		b |= Bit(1) << (f + r*8)
		if doesBlock(f, r, block) {
			break
		}
	}
	for r, f := rk-1, fl-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		b |= (Bit(1) << (f + r*8))
		if doesBlock(f, r, block) {
			break
		}
	}
	return b
}

// indexToBit will project the bits of index onto the given mask.
// So, if we have an mask of 0xF0F0, and an index of 0x1E, it will return 0x10E0.
func indexToBit(index int, mask Bit) (r Bit) {
	for i := 0; mask != 0; i++ {
		b := Bit(1) << bits.TrailingZeros64(uint64(mask))
		if index&(1<<i) != 0 {
			r |= b
		}
		mask ^= b
	}
	return r
}

func calcMagic(r *rand.Rand, sq int, maskF func(int) Bit, attF func(int, Bit) Bit) (Magic, [][]Coord, []Bit) {
	mask := maskF(sq)
	n := mask.CountOnes()
	a := make([]Bit, 1<<n)
	b := make([]Bit, 1<<n)
	used := make([]Bit, 1<<n)
	for i := range b {
		b[i] = indexToBit(i, mask)
		a[i] = attF(sq, b[i])
	}

	// Find the magic.
	var magic Bit
	for loop := true; loop; {
		magic = Bit(r.Uint64() & r.Uint64() & r.Uint64())
		if ((magic * mask) & (Bit(0xFF) << 56)).CountOnes() < 6 {
			continue
		}
		for i := range used {
			used[i] = 0
		}
		loop = false
		for i := range b {
			j := (int)((b[i] * magic) >> (64 - n))
			if used[j] == 0 {
				used[j] = a[i]
			} else if used[j] != a[i] {
				loop = true
				break
			}
		}
	}

	// With the magic, create the slice of attacked Coords.
	at := make([][]Coord, len(a))
	bits := make([]Bit, len(a))
	for k, v := range a {
		loc := (b[k] * magic) >> (64 - n)
		at[loc] = v.ToCoordSlice()
		bits[loc] = v
	}
	return Magic{Mask: mask, Value: magic, Shift: uint(64 - n)}, at, bits
}

func gen(w io.Writer) {
	genMagic := func(maskF func(int) Bit,
		attF func(int, Bit) Bit) (res [64]Magic, coords [64][][]Coord, bits [64][]Bit) {
		var wg sync.WaitGroup
		for i := 0; i < 64; i++ {
			wg.Add(1)
			go func(r *rand.Rand, j int) {
				res[j], coords[j], bits[j] = calcMagic(r, j, maskF, attF)
				wg.Done()
			}(rand.New(rand.NewSource(rand.Int63())), i)
		}
		wg.Wait()
		return res, coords, bits
	}
	rMagic, rCoords, rBits := genMagic(genRookMask, genRookAttacks)
	bMagic, bCoords, bBits := genMagic(genBishopMask, genBishopAttacks)

	writeMagic := func(name string, res [64]Magic) {
		fmt.Printf("size of %s: %d\n", name, getSize(res))
		fmt.Fprintf(w, "var %s = [64]Magic {\n", name)
		for i := 0; i < 64; i++ {
			fmt.Fprintf(w, "\tMagic { // %d\nMask: %v,\nValue: %v,\nShift: %d,\n},\n",
				i, res[i].Mask, res[i].Value, res[i].Shift)
		}
		fmt.Fprintf(w, "}\n\n")
	}
	writeCoords := func(name string, coords [64][][]Coord) {
		fmt.Printf("size of %s: %d\n", name, getSize(coords))
		fmt.Fprintf(w, "var %s = [64][][]Coord {\n", name)
		for i := 0; i < 64; i++ {
			fmt.Fprintf(w, "\t[][]Coord {\n")
			for j := range coords[i] {
				fmt.Fprintf(w, "\t\t[]Coord {")
				for k := range coords[i][j] {
					fmt.Fprintf(w, "%d, ", int(coords[i][j][k]))
				}
				fmt.Fprintf(w, "\t\t},\n")
			}
			fmt.Fprintf(w, "\t},\n")
		}
		fmt.Fprintf(w, "}\n\n")
	}
	writeBits := func(name string, bits [64][]Bit) {
		fmt.Printf("size of %s: %d\n", name, getSize(bits))
		fmt.Fprintf(w, "var %s = [64][]Bit {\n", name)
		for i := 0; i < 64; i++ {
			fmt.Fprintf(w, "\t[]Bit {\n")
			for _, v := range bits[i] {
				fmt.Fprintf(w, "%v,", v)
			}
			fmt.Fprintf(w, "\t},\n")
		}
		fmt.Fprintf(w, "}\n\n")
	}
	writeMagic("rookMagic", rMagic)
	writeCoords("rookCoords", rCoords)
	writeBits("rookBits", rBits)
	writeMagic("bishopMagic", bMagic)
	writeCoords("bishopCoords", bCoords)
	writeBits("bishopBits", bBits)
}

var header = `package main
// Code generated by go generate. DO NOT EDIT.

// rookLookup takes a coordinate and an occupancy bitset,
// returning a slice of coordinates we need to search for pieces.
func rookLookup(c Coord, occ Bit) []Coord {
	m := rookMagic[c.Idx()]
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	return rookCoords[c.Idx()][key]
}

func rookBit(c Coord, occ Bit) Bit {
	m := rookMagic[c.Idx()]
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	return rookBits[c.Idx()][key]
}

// bishopLookup takes a coordinate and an occupancy bitset,
// returning a slice of coordinates we need to search for pieces.
func bishopLookup(c Coord, occ Bit) []Coord {
	m := bishopMagic[c.Idx()]
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	return bishopCoords[c.Idx()][key]
}

func bishopBit(c Coord, occ Bit) Bit {
	m := bishopMagic[c.Idx()]
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	return bishopBits[c.Idx()][key]
}
`

func main() {
	b := bytes.NewBuffer([]byte(header))
	gen(b)

	out, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("magic_tables.go", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
