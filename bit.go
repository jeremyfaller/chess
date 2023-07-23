package main

import (
	"fmt"
	"math/bits"
	"strings"
)

type Bit uint64

// BitString pretty-prints a Bit.
func (b Bit) BitString() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%v\n", b.String()))
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

// String returns a formatted bitstring.
func (b Bit) String() string {
	return fmt.Sprintf("0x%016x", uint64(b))
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
