package main

//go:generate go run magic_gen.go magic.go bit.go coord.go dir.go

// Goal:
// magic := rookMagicBitboards[rookPos]
// blockers := occupied & magic.preMask
// cacheKey := (blockers * magic.magicValue) >> magic.rotate
// return magic.mapping[cacheKey]

// Magic
type Magic struct {
	Mask  Bit
	Value Bit
	Shift uint
}
