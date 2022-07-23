package main

import (
	"fmt"
	"strings"
)

type Piece uint8

const (
	Empty  Piece = 0
	Pawn         = 1
	Knight       = 2
	Bishop       = 3
	Rook         = 4
	Queen        = 5
	King         = 6
	White        = 8
	Black        = 16
)

type Score int

const (
	maxScore  = 10000000000
	minScore  = -maxScore
	stalemate = 0
	checkmate = minScore
)

var scores = []Score{
	0,       // Empty
	100,     // Pawn
	300,     // Knight
	300,     // Bishop
	500,     // Rook
	900,     // Queen
	100000,  // King
	0,       // Unused
	0,       // White|Empty
	100,     // White|Pawn
	300,     // White|Knight
	300,     // White|Bishop
	500,     // White|Rook
	900,     // White|Queen
	100000,  // White|King
	0,       // Unused
	0,       // Black|Empty
	-100,    // Black|Pawn
	-300,    // Black|Knight
	-300,    // Black|Bishop
	-500,    // Black|Rook
	-900,    // Black|Queen
	-100000, // Black|King
}

func (p Piece) IsWhite() bool {
	return p&White != 0
}

func (p Piece) IsBlack() bool {
	return p&Black != 0
}

func (p Piece) Colorless() Piece {
	return p &^ (White | Black)
}

func (p Piece) String() string {
	if p == Empty {
		return " "
	}
	printer := strings.ToUpper
	if p.IsBlack() {
		printer = strings.ToLower
	}
	switch p.Colorless() {
	case Pawn:
		return printer("P")
	case Knight:
		return printer("N")
	case Bishop:
		return printer("B")
	case Rook:
		return printer("R")
	case Queen:
		return printer("Q")
	case King:
		return printer("K")
	}
	if p.Color() == White {
		return "WHITE"
	}
	return "BLACK"
}

func (p Piece) IsEmpty() bool {
	return p == Empty
}

func (p Piece) IsKing() bool {
	return p.Colorless() == King
}

func (p Piece) IsRook() bool {
	return p.Colorless() == Rook
}

func (p Piece) IsPawn() bool {
	return p.Colorless() == Pawn
}

func (p Piece) Score() Score {
	return scores[p]
}

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

// Attack dir returns a slice of Dir in which a Piece attacks.
func (p Piece) AttackDir() []Dir {
	switch p.Colorless() {
	case Queen, Rook, Bishop, Knight:
		return p.MoveDir()
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

// MoveDir returns a slice of Dir in which a Piece moves.
func (p Piece) MoveDir() []Dir {
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

// isSlider returns true if a piece is a sliding piece, ie it can move more
// than one space in a given direciton.
func (p Piece) isSlider() bool {
	switch p.Colorless() {
	case Bishop, Rook, Queen:
		return true
	default:
		return false
	}
}

// AttackDistance returns the distance a piece can attack in a given direction.
func (p Piece) AttackDistance(d Dir) int {
	switch p.Colorless() {
	case Knight, King, Pawn:
		return 1
	default:
		return 8
	}
}

// SlideDistance returns the distance a piece can slide.
func (p Piece) SlideDistance() int {
	if !p.isSlider() {
		return 1
	}
	return 8 // can slide upto the whole board.
}

func (p Piece) NoteString() string {
	switch p.Colorless() {
	case Pawn:
		return ""
	case King:
		return "K"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Queen:
		return "Q"
	case Rook:
		return "R"
	}
	panic(fmt.Sprintf("unknown piece %d", p))
}

func (p Piece) Color() Piece {
	return p & (White | Black)
}

func (p Piece) OppositeColor() Piece {
	return p.Color() ^ (White | Black)
}

// HashIdx returns number [0..11] for the piece.
func (p Piece) HashIdx() int {
	if p == Empty {
		panic("empty hash")
	}
	c := p.Colorless() - Pawn
	if p.Color() == Black {
		c += 6
	}
	return int(c)
}

func (s Score) String() string {
	return fmt.Sprintf("%01.2f", float32(s)/100.)
}
