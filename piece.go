package main

import (
	"fmt"
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
	White        = 0
	Black        = 8
)

//go:generate go run piece.go piece_gen.go dir.go coord.go bit.go

type Score int16

var scores = []Score{
	0,      // White|Empty
	100,    // White|Pawn
	300,    // White|Knight
	300,    // White|Bishop
	500,    // White|Rook
	900,    // White|Queen
	10000,  // White|King
	0,      // Unused
	0,      // Black|Empty
	-100,   // Black|Pawn
	-300,   // Black|Knight
	-300,   // Black|Bishop
	-500,   // Black|Rook
	-900,   // Black|Queen
	-10000, // Black|King
}

var runeToColorlessPiece = map[rune]Piece{
	'p': Pawn,
	'n': Knight,
	'b': Bishop,
	'r': Rook,
	'q': Queen,
	'k': King,
}

func (p Piece) IsWhite() bool {
	return p&Black == 0
}

func (p Piece) IsBlack() bool {
	return p&Black != 0
}

func (p Piece) Colorless() Piece {
	return p &^ Black
}

func (p Piece) String() string {
	//black, white, reset := "\033[31m", "\033[37m", "\033[0m"
	black, white, reset := "", "", ""
	switch p {
	case White | Pawn:
		return white + "P" + reset
	case White | Knight:
		return white + "N" + reset
	case White | Bishop:
		return white + "B" + reset
	case White | Rook:
		return white + "R" + reset
	case White | Queen:
		return white + "Q" + reset
	case White | King:
		return white + "K" + reset
	case Black | Pawn:
		return black + "p" + reset
	case Black | Knight:
		return black + "n" + reset
	case Black | Bishop:
		return black + "b" + reset
	case Black | Rook:
		return black + "r" + reset
	case Black | Queen:
		return black + "q" + reset
	case Black | King:
		return black + "k" + reset
	}
	return " "
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

func (p Piece) NoteString() string {
	switch p.Colorless() {
	case Pawn:
		return "p"
	case King:
		return "k"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Queen:
		return "q"
	case Rook:
		return "r"
	}
	return "?"
}

func (p Piece) Color() Piece {
	return p & Black
}

func (p Piece) OppositeColor() Piece {
	return (p & Black) ^ Black
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

// HashIdx returns number [0..11] for the piece.
func (p Piece) HashIdx() int {
	if p == Empty {
		panic("empty hash")
	}
	c := p.Colorless() - Pawn
	if p.Color() == Black {
		c += 6
	}
	return int(c) * 64
}

func (s Score) String() string {
	return fmt.Sprintf("%01.2f", float32(s)/100.)
}
