package main

import (
	"fmt"
	"math/bits"
	"strings"
)

type Bit uint64

// String pretty-prints a Bit.
func (b Bit) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%016x\n", b.Uint64()))
	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			if b&(1<<(x+y*8)) == 0 {
				s.WriteByte('0')
			} else {
				s.WriteByte('1')
			}
			s.WriteByte(' ')
		}
		if y != 0 {
			s.WriteByte('\n')
		}
	}
	return s.String()
}

// Uint64 returns the uint64 value for this Bit.
func (b Bit) Uint64() uint64 {
	return uint64(b)
}

// IsSet returns true is if bit at the given index is set.
func (b Bit) IsSet(idx int) bool {
	return b&(1<<idx) != 0
}

// Set sets a given bit index.
func (b *Bit) Set(idx int) Bit {
	*b |= Bit(1 << idx)
	return *b
}

// Clear clears a given bit index.
func (b *Bit) Clear(idx int) Bit {
	*b &= Bit(^(1 << idx))
	return *b
}

// And ands two Bit.
func (b *Bit) And(b2 Bit) Bit {
	*b &= b2
	return *b
}

// CountOnes returns the number of ones.
func (b Bit) CountOnes() int {
	return bits.OnesCount64(uint64(b))
}
