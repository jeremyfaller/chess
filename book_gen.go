//go:build ignore
// +build ignore

package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

//go:embed book.txt
var book string

var header = `package main

import (
	"math/rand"
)

func getBook(b *Board, r *rand.Rand) (Move, bool) {
	if v, ok := book[b.ZHash()]; ok {
		idx := r.Intn(v.total)
		for _, test := range v.moves {
			idx -= test.count
			if idx < 0 {
				move := Move{from:test.from, to:test.to, p:b.at(test.from)}
				if !b.isLegalMove(&move) {
					panic("move should be legal")
				}
				return move, true
			}
		}
		panic("shouldn't be reachable")
	}
	return Move{}, false
}

type shortMove struct {
	count 		int
	from, to	Coord
}

type Book struct {
	total	int
	moves	[]shortMove
}

`

func gen(f io.Writer) {
	lut := make(map[string]string)
	var lastPos string
	var moves []string
	save := func() {
		if _, ok := lut[lastPos]; ok {
			panic("unexpected")
		}
		if len(lastPos) != 0 {
			lut[lastPos] = strings.Join(moves, " ")
		}
		lastPos = ""
		moves = moves[:0]
	}

	// Read the file.
	for _, str := range strings.Split(book, "\n") {
		str = strings.Trim(str, " \t")
		if strings.HasPrefix(str, "pos ") {
			save()
			lastPos, _ = strings.CutPrefix(str, "pos ")
		} else if len(str) != 0 {
			moves = append(moves, str)
		}
	}
	save()

	b := New()
	seen := make(map[Hash]struct{})
	printMoves := func() bool {
		fen := strings.Join(strings.Split(b.FENString(), " ")[0:4], " ")
		if v, ok := lut[fen]; ok {
			seen[b.ZHash()] = struct{}{}
			fmt.Fprintf(f, "\t%d: Book{\n\t\tmoves: []shortMove{\n", b.ZHash())
			fields := strings.Split(v, " ")
			var total int
			for i := 0; i < len(fields); i += 2 {
				mStr, cStr := fields[i], fields[i+1]
				m, err := b.parseAlgebraic(mStr)
				if err != nil {
					panic(err)
				}
				c, err := strconv.Atoi(cStr)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(f, "\t\t\t{ from: %d, to: %d, count: %d },\n", m.from.Idx(), m.to.Idx(), c)
				total += c
			}
			fmt.Fprintf(f, "\t\t},\n\t\ttotal: %d,\n", total)
			fmt.Fprintf(f, "\t},\n")
			return true
		}
		return false
	}

	// Now parse the moves.
	fmt.Fprintf(f, "var book = map[Hash]Book {\n")
	printMoves()
	var deep func()
	deep = func() {
		for _, move := range b.PossibleMoves(nil) {
			b.MakeMove(move)
			if _, ok := seen[b.ZHash()]; ok {
				b.UnmakeMove()
				continue
			}
			if printMoves() {
				deep()
			}
			b.UnmakeMove()
		}
	}
	deep()
	fmt.Fprintf(f, "}\n\n")
}

func main() {
	flag.Parse()
	b := bytes.NewBuffer([]byte(header))
	gen(b)

	out, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("book_tables.go", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
