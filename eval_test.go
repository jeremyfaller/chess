package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

//go:embed testdata/mate.dat
var mate string

func TestMateIn(t *testing.T) {
	// Read the test data.
	type testT struct {
		depth int
		fen   string
	}
	var tests []testT
	for _, str := range strings.Split(mate, "\n") {
		if len(str) == 0 {
			continue
		}
		depth, err := strconv.Atoi(str[0:1])
		if err != nil {
			t.Fatal(err)
		}
		tests = append(tests, testT{depth: 2*depth - 1, fen: str[2:]})
	}

	// Run the tests.
	for i, test := range tests {
		test := test

		t.Run(fmt.Sprintf("test %d, mate in %d", i, test.depth), func(t *testing.T) {
			t.Parallel()
			b, err := FromFEN(test.fen)
			if err != nil {
				t.Fatalf("[%d] error on fen %v", i, err)
			}

			e := NewEval(b, test.depth)
			e.Start()
			if e.score != -checkmate {
				t.Errorf("[%d] was not a checkmate %v", i, test.fen)
			}
		})
	}
}
