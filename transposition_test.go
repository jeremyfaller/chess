package main

import (
	"testing"
)

func TestTranspositionSize(t *testing.T) {
	t.SkipNow()
	tt := NewTranspositionTable(1)
	s1 := tt.Size()
	tt.Resize(10)
	s2 := tt.Size()
	t.Logf("elements before: %d, elements after: %d", s1, s2)
	if s1 == 0 || s2 == 0 || s1 >= s2 || s1*10 > s2 {
		t.Errorf("table didn't Resize properly s1: %d, s2: %d", s1, s2)
	}
}
