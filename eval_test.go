package main

import (
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		desc string
		fen  string
	}{
		{"mate in 2", "r2qkb1r/pp2nppp/3p4/2pNN1B1/2BnP3/3P4/PPP2PPP/R2bK2R w KQkq - 0 1"},
	}

	for i, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%d] %s error on fen %v", i, test.desc, err)
		}
		e := NewEval(b)
		e.Start()
		t.Logf("%+v\n", e)
		t.Logf("perft(3) = %d", b.Perft(3))
	}
}
