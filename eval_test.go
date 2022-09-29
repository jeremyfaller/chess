package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

//go:embed testdata/mate.dat
var mates string

type evalTest struct {
	depth int
	fen   string
}

// getTests returns testTs from the passed in string.
func getTests(dat string) []evalTest {
	var tests []evalTest
	for _, str := range strings.Split(dat, "\n") {
		if len(str) == 0 {
			continue
		}
		depth, err := strconv.Atoi(str[0:1])
		if err != nil {
			panic(err)
		}
		tests = append(tests, evalTest{depth: 2*depth - 1, fen: str[2:]})
	}
	return tests
}

func TestMateIn(t *testing.T) {
	// Run the tests.
	for i, test := range getTests(mates) {
		test := test

		t.Run(fmt.Sprintf("test %d, mate in %d", i, test.depth), func(t *testing.T) {
			t.Parallel()
			b, err := FromFEN(test.fen)
			if err != nil {
				t.Fatalf("[%d] error in fen %v", i, err)
			}

			e := NewEval(b, test.depth)
			e.Start()
			if e.score != -checkmate {
				t.Errorf("[%d] was not a checkmate %v", i, test.fen)
			}
		})
	}
}

func mateBenchmarker(b *testing.B, d int, tests []evalTest) {
	for j := 0; j < b.N; j++ {
		for i, test := range getTests(mates) {
			if test.depth != d {
				continue
			}
			b, err := FromFEN(test.fen)
			if err != nil {
				panic(fmt.Sprintf("[%d] error in fen %v", i, err))
			}
			e := NewEval(b, test.depth)
			e.Start()
			if e.score != -checkmate {
				panic(fmt.Sprintf("[%d] was not a checkmate %v", i, test.fen))
			}
		}
	}
}

func BenchmarkMateIn2(b *testing.B) {
	mateBenchmarker(b, 2*2-1, getTests(mates))
}

func BenchmarkMateIn3(b *testing.B) {
	mateBenchmarker(b, 3*2-1, getTests(mates))
}
