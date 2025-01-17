// Generates Piece movement.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
)

var (
	kingAttackDir      = []Dir{N, NE, E, SE, S, SW, W, NW}
	kingDir            = []Dir{N, NE, E, E2, SE, S, SW, W, W2, NW}
	queenDir           = []Dir{N, NE, E, SE, S, SW, W, NW}
	rookDir            = []Dir{N, E, S, W}
	bishopDir          = []Dir{NE, SE, SW, NW}
	knightDir          = []Dir{NNE, NEE, SEE, SSE, SSW, SWW, NWW, NNW}
	whitePawnDir       = []Dir{N, NN, NE, NW}
	whitePawnAttackDir = []Dir{NE, NW}
	blackPawnDir       = []Dir{S, SS, SE, SW}
	blackPawnAttackDir = []Dir{SE, SW}
)

// attackDir returns a slice of Dir in which a Piece attacks.
func (p Piece) attackDir() []Dir {
	switch p.Colorless() {
	case Queen, Rook, Bishop, Knight:
		return p.moveDir()
	case King:
		return kingAttackDir
	case Pawn:
		if p.Color() == White {
			return whitePawnAttackDir
		} else {
			return blackPawnAttackDir
		}
	}
	panic("direction not set up for piece " + p.String())
}

// moveDir returns a slice of Dir in which a Piece moves.
func (p Piece) moveDir() []Dir {
	switch p.Colorless() {
	case Queen:
		return queenDir
	case King:
		return kingDir
	case Rook:
		return rookDir
	case Bishop:
		return bishopDir
	case Pawn:
		if p.Color() == White {
			return whitePawnDir
		} else {
			return blackPawnDir
		}
	case Knight:
		return knightDir
	}
	panic("direction not set up for piece " + p.String())
}

// attackDistance returns the distance a piece can attack in a given direction.
func (p Piece) attackDistance(d Dir) int {
	switch p.Colorless() {
	case Knight, King, Pawn:
		return 1
	default:
		return 8
	}
}

// genMoves generates the set of moves for a piece at a coordinate.
func genMoves(p Piece, c Coord) (moves []Coord) {
	if p.isSlider() {
		return moves
	}
	for _, d := range p.moveDir() {
		pos := c.ApplyDir(d)
		if !pos.IsValid() {
			continue
		}
		moves = append(moves, pos)
	}
	return moves
}

// genAttacks returns a bit for a piece and a location.
func genAttacks(p Piece, c Coord) (b Bit) {
	for _, d := range p.attackDir() {
		pos := c.ApplyDir(d)
		if !pos.IsValid() {
			continue
		}
		b.Set(pos.Idx())
	}
	return b
}

func gen(w io.Writer) {
	// Make the moves LUT.
	movesForPiece := make([][64][]Coord, Black*2)
	for _, c := range []Piece{White, Black} {
		for _, p := range []Piece{Pawn, Knight, King} {
			p |= c
			for i := 0; i < 64; i++ {
				coord := CoordFromIdx(i)
				movesForPiece[p][i] = genMoves(p, coord)
			}
		}
	}

	fmt.Fprintf(w, "var movesForPiece = [][64][]Coord {\n")
	for i := range movesForPiece {
		fmt.Fprintf(w, "\t[64][]Coord {\n")
		p := Piece(i).Colorless()
		if p == Pawn || p == Knight || p == King {
			for j := range movesForPiece[i] {
				fmt.Fprintf(w, "\t\t[]Coord {")
				for k := range movesForPiece[i][j] {
					fmt.Fprintf(w, " %d,", movesForPiece[i][j][k].Idx())
				}
				fmt.Fprintf(w, "\t\t},\n")
			}
		}
		fmt.Fprintf(w, "\t},\n")
	}
	fmt.Fprintf(w, "}\n\n")

	writeAttacks := func(name string, p Piece) {
		fmt.Fprintf(w, "var %sAttacks = [64]Bit {", name)
		for i := 0; i < 64; i++ {
			fmt.Fprintf(w, "%s, ", genAttacks(p, CoordFromIdx(i)))
		}
		fmt.Fprintf(w, "}\n\n")
	}
	writeAttacks("wPawn", Pawn|White)
	writeAttacks("bPawn", Pawn|Black)
	writeAttacks("king", King)
	writeAttacks("knight", Knight)
}

func main() {
	b := bytes.NewBuffer([]byte(header))
	gen(b)

	out, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("piece_tables.go", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

var header = `package main
// Code generated by go generate. DO NOT EDIT.

// Moves returns a slice of slices of all the squares a piece could possibly move for
// a Piece at a given Coord.
func (p Piece) Moves(c Coord, occ Bit) []Coord {
	if p.isSlider() {
		if p.Colorless() == Bishop {
			return bishopLookup(c, occ)
		}
		return rookLookup(c, occ)
	}
	return movesForPiece[p][c.Idx()]
}

// Attacks returns a Bit of the attacked squares for a given piece.
func (p Piece) Attacks(c Coord, occ Bit) Bit {
	idx := c.Idx()
	switch p {
		case White|Pawn:
			return wPawnAttacks[idx]
		case Black|Pawn:
			return bPawnAttacks[idx]
		case White|Knight, Black|Knight:
			return knightAttacks[idx]
		case White|Bishop, Black|Bishop:
			return bishopBit(c, occ)
		case White|Rook, Black|Rook:
			return rookBit(c, occ)
		case White|Queen, Black|Queen:
			return rookBit(c, occ) | bishopBit(c, occ)
		case White|King, Black|King:
			return kingAttacks[idx]
	}
	panic("unknown piece")
}
`
