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
	"reflect"
	"sync"
)

func getSize(v interface{}) int {
	size := int(reflect.TypeOf(v).Size())
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			size += getSize(s.Index(i).Interface())
		}
	case reflect.Map:
		s := reflect.ValueOf(v)
		keys := s.MapKeys()
		size += int(float64(len(keys)) * 10.79) // approximation from https://golang.org/src/runtime/hashmap.go
		for i := range keys {
			size += getSize(keys[i].Interface()) + getSize(s.MapIndex(keys[i]).Interface())
		}
	case reflect.String:
		size += reflect.ValueOf(v).Len()
	case reflect.Struct:
		s := reflect.ValueOf(v)
		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).CanInterface() {
				size += getSize(s.Field(i).Interface())
			}
		}
	case reflect.Int, reflect.Uint64:
		size += 8
	default:
		panic(fmt.Sprintf("UNKNOWN %v", reflect.TypeOf(v).Kind()))
	}
	return size
}

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

func calcMagic(r *rand.Rand, sq int, maskF func(int) Bit, attF func(int, Bit) Bit) (Magic, [][]Coord) {
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
	for k, v := range a {
		loc := (b[k] * magic) >> (64 - n)
		at[loc] = make([]Coord, v.CountOnes())
		for j := range at[loc] {
			n := bits.TrailingZeros64(uint64(v))
			at[loc][j] = CoordFromIdx(n)
			v ^= Bit(1 << n)
		}
	}
	return Magic{Mask: mask, Value: magic, Shift: (64 - n)}, at
}

func gen(w io.Writer) {
	genMagic := func(maskF func(int) Bit,
		attF func(int, Bit) Bit) (res [64]Magic, coords [64][][]Coord) {
		var wg sync.WaitGroup
		for i := 0; i < 64; i++ {
			wg.Add(1)
			go func(r *rand.Rand, j int) {
				res[j], coords[j] = calcMagic(r, j, maskF, attF)
				wg.Done()
			}(rand.New(rand.NewSource(rand.Int63())), i)
		}
		wg.Wait()
		return res, coords
	}
	rMagic, rCoords := genMagic(genRookMask, genRookAttacks)
	bMagic, bCoords := genMagic(genBishopMask, genBishopAttacks)

	writeMagic := func(name string, res [64]Magic) {
		fmt.Printf("size of %s: %d\n", name, getSize(res))
		fmt.Fprintf(w, "var %s = [64]Magic {\n", name)
		for i := 0; i < 64; i++ {
			fmt.Fprintf(w, "\tMagic { // %d\nMask: 0x%016x,\nValue: 0x%016x,\nShift: %d,\n},\n",
				i, res[i].Mask.Uint64(), res[i].Value.Uint64(), res[i].Shift)
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
	writeMagic("rookMagic", rMagic)
	writeCoords("rookCoords", rCoords)
	writeMagic("bishopMagic", bMagic)
	writeCoords("bishopCoords", bCoords)
}

var header = `package main
// Code generated by go generate. DO NOT EDIT.

// RookLookup takes a coordinate and an occupancy bitset,
// returning a slice of coordinates we need to search for pieces.
func RookLookup(c Coord, occ Bit) []Coord {
	if occ&(1<<c.Idx()) == 0 {
		panic("expected rook to be at that space")
	}
	m := rookMagic[c.Idx()]
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	return rookCoords[c.Idx()][key]
}

// BishopLookup takes a coordinate and an occupancy bitset,
// returning a slice of coordinates we need to search for pieces.
func BishopLookup(c Coord, occ Bit) []Coord {
	if occ&(1<<c.Idx()) == 0 {
		panic("expected bishop to be at that space")
	}
	m := bishopMagic[c.Idx()]
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	return bishopCoords[c.Idx()][key]
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
