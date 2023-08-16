package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

//go:embed testdata/mate.dat
var mates string

type evalTest struct {
	depth Depth
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
		tests = append(tests, evalTest{depth: Depth(2*depth - 1), fen: str[2:]})
	}
	return tests
}

func TestMateIn(t *testing.T) {
	queue := make(chan struct{}, 10)

	// Run the tests.
	for i, test := range getTests(mates) {
		test := test
		name := fmt.Sprintf("test %d, mate in %d", i, test.depth)

		queue <- struct{}{}
		t.Run(name, func(t *testing.T) {
			<-queue
			t.Parallel()

			if test.depth >= 7 {
				t.Skip("skipping: " + t.Name() + " because it's too long")
			}

			b, err := FromFEN(test.fen)
			if err != nil {
				t.Fatalf("[%d] error in fen %v", i, err)
			}

			e := NewEval(test.depth)
			e.Start(b)
			e.Wait()
			if !IsMateScore(e.score) {
				t.Errorf("[%d] was not a checkmate %v, %v", i, test.fen, e.score)
			}
		})
	}
}

func TestIsMateScore(t *testing.T) {
	tests := []struct {
		s      Score
		isMate bool
	}{
		{checkmate, true},
		{-checkmate, true},
		{0, false},
	}

	for i, test := range tests {
		if v := IsMateScore(test.s); v != test.isMate {
			t.Errorf("[%d] IsMateScore(%v) = %t, expected %v", i, test.s, v, test.isMate)
		}
	}
}

func TestEvalCancel(t *testing.T) {
	e := NewEval(100)
	e.Start(New())
	if !e.IsRunning() {
		t.Fatalf("expected eval running")
	}
	e.Stop()
	if e.IsRunning() {
		t.Fatalf("expected eval stopped")
	}
}

func TestEvalTimeout(t *testing.T) {
	t.Parallel()
	dur := 10 * time.Millisecond
	e := NewEval(100)
	e.SetDuration(dur)
	e.Start(New())
	if !e.IsRunning() {
		t.Fatalf("expected eval running")
	}
	time.Sleep(dur * 2)
	if e.IsRunning() {
		t.Fatalf("expected eval stopped")
	}
}

func mateBenchmarker(b *testing.B, d Depth, tests []evalTest) {
	for j := 0; j < b.N; j++ {
		for i, test := range getTests(mates) {
			if test.depth != d {
				continue
			}
			b, err := FromFEN(test.fen)
			if err != nil {
				panic(fmt.Sprintf("[%d] error in fen %v", i, err))
			}
			e := NewEval(test.depth)
			e.Start(b)
			e.Wait()
			if !IsMateScore(e.score) {
				panic(fmt.Sprintf("[%d] was not a checkmate %v, %d", i, test.fen, e.score))
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

func BenchmarkMateIn4(b *testing.B) {
	mateBenchmarker(b, 4*2-1, getTests(mates))
}
