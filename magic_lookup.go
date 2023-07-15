package main

import "fmt"

// RookLookup takes a coordinate and an occupancy bitset,
// returning a slice of coordinates we need to search for pieces.
func RookLookup(c Coord, occ Bit) []Coord {
	if occ&(1<<c.Idx()) == 0 {
		panic("expected rook to be at that space")
	}
	m := rookMagic[c.Idx()]
	fmt.Println("occ\n", occ)
	fmt.Println("mask\n", m.Mask)
	fmt.Println("and\n", m.Mask&occ)
	key := ((m.Mask & occ) * m.Value) >> m.Shift
	fmt.Println("key", key.Uint64())
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
