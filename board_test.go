package main

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func makeMove(t *testing.T, desc, before, from, to string) (*Board, Move) {
	t.Helper()
	coord := testingCoordFunc(t)

	b, err := FromFEN(before)
	if err != nil {
		t.Errorf("[%v] error creating board", err)
	}
	if diff := pretty.Compare(before, b.String()); diff != "" {
		t.Errorf("[%s] boards unequal\n%s", desc, diff)
	}

	m := Move{from: coord(from), to: coord(to), p: b.at(coord(from))}
	if !b.isLegalMove(&m) {
		t.Errorf("[%v] isLegalMove == false", desc)
	}
	b.MakeMove(m)

	return b, m
}

func makeUnmakeMove(t *testing.T, desc, before, from, to, after string) (*Board, Move) {
	b, m := makeMove(t, desc, before, from, to)

	if len(after) != 0 {
		if diff := pretty.Compare(after, b.String()); diff != "" {
			t.Errorf("[%s] after boards unequal\n%s", desc, diff)
		}
	}

	b.UnmakeMove()
	if diff := pretty.Compare(before, b.String()); diff != "" {
		t.Errorf("[%s] unmake boards unequal\n%s", desc, diff)
	}

	return b, m
}

func TestIsLegalMove(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc      string
		fen       string
		move      Move
		isLegal   bool
		isCapture bool
	}{
		{ // Check that pawn captures are valid with a piece there.
			"pawn capture",
			"k7/8/8/p7/1P6/8/8/K7 w - - 0 1",
			Move{p: White | Pawn, to: coord("a5"), from: coord("b4")},
			true,
			true,
		},
		{ // Check that diagonal moves are invalid when there's no piece.
			"pawn cannot capture",
			"k7/8/8/8/1P6/8/8/K7 w - - 0 1",
			Move{p: White | Pawn, to: coord("a5"), from: coord("b4")},
			false,
			false,
		},
		{ // Check that a pawn can't move through a piece.
			"pawn cannot move through piece",
			"k7/8/8/8/8/N7/P7/K7 w - - 0 1",
			Move{p: White | Pawn, to: coord("a4"), from: coord("a2")},
			false,
			false,
		},
		{ // Check that en passant moves register as valid.
			"en passant",
			"k7/8/8/pP6/8/8/8/K7 w - a6 0 1",
			Move{p: White | Pawn, to: coord("a6"), from: coord("b5")},
			true,
			true,
		},
		{ // Check that en passant moves are invalid when the enpassant target isn't set.
			"non-en passant",
			"k7/8/8/pP6/8/8/8/K7 w - - 0 1",
			Move{p: White | Pawn, to: coord("a6"), from: coord("b5")},
			false,
			false,
		},
		{ // Check that you can't take pieces of your own color.
			"same color is illegal capture",
			"Rk6/8/8/8/8/8/8/RK6 w - - 0 1",
			Move{p: White | Rook, to: coord("a8"), from: coord("a1")},
			false,
			false,
		},
		{ // Check that you can't castle out of check.
			"castling from check is illegal",
			"k7/8/8/8/8/8/4r3/R3K2R w KQ - 0 1",
			Move{p: White | King, to: coord("g1"), from: coord("e1")},
			false,
			false,
		},
		{ // Can't castle across a square in check.
			"illegal kingside castle",
			"k7/8/8/8/8/5r2/8/R3K2R w KQ - 0 1",
			Move{p: White | King, to: coord("g1"), from: coord("e1")},
			false,
			false,
		},
		{ // Can't move into check.
			"can't move into check",
			"k7/8/8/8/8/5r2/8/R3K2R w KQ - 0 1",
			Move{p: White | King, to: coord("f1"), from: coord("e1")},
			false,
			false,
		},
		{ // Can't move a piece that would result in check.
			"move results in check",
			"k7/8/8/8/8/4r3/4R3/4K3 w - - 0 1",
			Move{p: White | Rook, to: coord("a2"), from: coord("e2")},
			false,
			false,
		},
	}

	for _, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Errorf("[%s] error creating board: %v", test.desc, err)
		}
		if isLegal := b.isLegalMove(&test.move); isLegal != test.isLegal {
			t.Errorf("[%s] isLegalMove = %t, expected %t", test.desc, isLegal, test.isLegal)
		}
		if test.move.isCapture != test.isCapture {
			t.Errorf("[%s] move.isCapture = %t, expected %t", test.desc, test.move.isCapture, test.isCapture)
		}
	}
}

func TestAddMoves(t *testing.T) {
	coord := testingCoordFunc(t)
	move := func(p Piece, from, to string, isCap bool) Move {
		return Move{p: p, from: coord(from), to: coord(to), isCapture: isCap}
	}
	promo := func(p Piece, from, to string, promo Piece) Move {
		return Move{p: p, from: coord(from), to: coord(to), promotion: promo}
	}
	toSet := func(moves []Move) map[Move]struct{} {
		mapSet := make(map[Move]struct{})
		for _, m := range moves {
			mapSet[m] = struct{}{}
		}
		return mapSet
	}
	tests := []struct {
		desc  string
		fen   string
		c     Coord
		moves map[Move]struct{}
	}{
		{ // Test that an opening pawn move includes 1&2 space moves.
			"opening pawn moves",
			"k7/8/8/8/8/8/P7/K w - - 0 1",
			coord("a2"),
			toSet([]Move{
				move(White|Pawn, "a2", "a4", false),
				move(White|Pawn, "a2", "a3", false),
			}),
		},
		{ // Test that a non-opening pawn move is only one space.
			"pawn move",
			"k7/8/8/8/8/P7/8/K7 w - - 0 1",
			coord("a3"),
			toSet([]Move{
				move(White|Pawn, "a3", "a4", false),
			}),
		},
		{ // Test moving into your own piece stops the move.
			"same color stops",
			"k7/8/8/8/8/R7/8/R1K5 w - - 0 1",
			coord("a1"),
			toSet([]Move{
				move(White|Rook, "a1", "a2", false),
				move(White|Rook, "a1", "b1", false),
			}),
		},
		{
			"king moves",
			"k7/8/8/8/8/8/8/R3K2R w - - 0 1",
			coord("e1"),
			toSet([]Move{
				move(White|King, "e1", "d1", false),
				move(White|King, "e1", "d2", false),
				move(White|King, "e1", "e2", false),
				move(White|King, "e1", "f1", false),
				move(White|King, "e1", "f2", false),
				move(White|King, "e1", "g1", false),
				move(White|King, "e1", "c1", false),
			}),
		},
		{
			"promotions",
			"k7/4P3/8/8/8/8/8/K7 w - - 0 1",
			coord("e7"),
			toSet([]Move{
				promo(White|Pawn, "e7", "e8", White|Queen),
				promo(White|Pawn, "e7", "e8", White|Knight),
				promo(White|Pawn, "e7", "e8", White|Bishop),
				promo(White|Pawn, "e7", "e8", White|Rook),
			}),
		},
		{
			"pawn cannot move into another piece",
			"rnbqkbnr/1ppppppp/8/p7/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2",
			coord("a4"),
			toSet([]Move{}),
		},
		{
			"king can't castle across piece",
			"rnbqkbnr/1ppppppp/p7/8/8/7N/PPPPPPPP/RNBQKB1R w KQkq - 0 2",
			coord("e1"),
			toSet([]Move{}),
		},
	}

	for _, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%s] error creating board; %v", test.desc, err)
		}
		moves := toSet(b.addMoves(test.c, make([]Move, 0)))
		if diff := pretty.Compare(test.moves, moves); diff != "" {
			t.Errorf("[%s] moves unequal:\n%s", test.desc, diff)
		}
	}
}

func TestFENChecks(t *testing.T) {
	tests := []struct {
		desc               string
		fen                string
		isWCheck, isBCheck bool
	}{
		{
			"no check",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			false,
			false,
		},
		{
			"white check",
			"1k6/8/8/8/8/4r3/8/R3K2R w - - 0 1",
			true,
			false,
		},
		{
			"black check",
			"r3k2r/8/4R3/8/8/8/8/1K6 b - - 0 1",
			false,
			true,
		},
	}

	for _, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%s] unexpected error: %v", test.desc, err)
		}
		if v := b.IsKingInCheck(White | King); v != test.isWCheck {
			t.Errorf("[%s] IsKingInCheck(White|King) = %v, expected = %v", test.desc, v, test.isWCheck)
		}
		if v := b.IsKingInCheck(Black | King); v != test.isBCheck {
			t.Errorf("[%s] IsKingInCheck(Black|King) = %v, expected = %v", test.desc, v, test.isBCheck)
		}
	}
}

func TestUpdateCastleState(t *testing.T) {
	tests := []struct {
		desc     string
		fen      string
		from, to string
		state    string
	}{
		{
			"a1xa8",
			"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			"a1", "a8",
			"Kk", // kingside rights retained.
		},
		{
			"h8xh1",
			"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			"h8", "h1",
			"Qq", // queenside rights retained.
		},
		{
			"Ke2",
			"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			"e1", "e2",
			"kq", // black retains rights.
		},
		{
			"Ke7",
			"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			"e8", "e7",
			"KQ", // white retains rights.
		},
	}

	for _, test := range tests {
		b, _ := makeMove(t, test.desc, test.fen, test.from, test.to)
		if s := b.castleString(); s != test.state {
			t.Errorf("[%s] expected %q to contain %q", test.desc, s, test.state)
		}
	}
}

func TestMakeMove(t *testing.T) {
	tests := []struct {
		desc     string
		before   string
		from, to string
		after    string
	}{
		{
			"e4",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"e2",
			"e4",
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			"c5",
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			"c7",
			"c5",
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		},
		{
			"Nf3",
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			"g1",
			"f3",
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		},
		{
			"kingside castle",
			"k7/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			"e1",
			"g1",
			"k7/8/8/8/8/8/8/R4RK1 b kq - 1 1",
		},
		{
			"queenside castle",
			"k7/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			"e1",
			"c1",
			"k7/8/8/8/8/8/8/2KR3R b kq - 1 1",
		},
		{
			"en passant capture",
			"k7/8/8/pP6/8/8/8/K7 w - a6 0 1",
			"b5",
			"a6",
			"k7/8/P7/8/8/8/8/K7 b - - 0 1",
		},
	}

	for _, test := range tests {
		makeUnmakeMove(t, test.desc, test.before, test.from, test.to, test.after)
	}
}

func TestAttackers(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc string
		p    Piece
		c    Coord
		ex   []string
	}{
		{
			"rook",
			White | Rook,
			coord("a1"),
			[]string{"a2", "a3", "a4", "a5", "a6", "a7", "a8", "b1", "c1", "d1", "e1", "f1", "g1", "h1"},
		},
		{
			"pawn",
			White | Pawn,
			coord("b2"),
			[]string{"a3", "c3"},
		},
	}

	for _, test := range tests {
		var a PsuedoMoves
		a.Update(test.p, test.c)
		attacked := make(map[Coord]struct{})
		for _, s := range test.ex {
			attacked[coord(s)] = struct{}{}
		}
		for i := 0; i < 64; i++ {
			c := CoordFromIdx(i)
			if at := a.Attackers(c); len(at) == 1 && at[0] == test.c {
				delete(attacked, c)
			}
		}
		if len(attacked) != 0 {
			t.Errorf("Expected to clear the attackers list: %v", attacked)
		}
	}
}

func TestDoesSquareAttack(t *testing.T) {
	coord := testingCoordFunc(t)
	tests := []struct {
		desc     string
		fen      string
		from, to Coord
		ex       bool
	}{
		{
			"rook to rook",
			"k7/8/8/8/8/r7/R7/K7 w - - 0 1",
			coord("a3"),
			coord("a2"),
			true,
		},
		{
			"rook to king",
			"k7/8/8/8/8/r7/R7/K7 w - - 0 1",
			coord("a3"),
			coord("a1"),
			false,
		},
	}

	for _, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%s] unexpected error: %v", test.desc, err)
		}
		if v := b.doesSquareAttack(test.from, test.to); v != test.ex {
			t.Errorf("[%s] doesSquareAttack(%v, %v) = %v, expected = %v", test.desc, test.from, test.to, v, test.ex)
		}
	}
}

func TestInvalidCoordNotAttacked(t *testing.T) {
	b := New()
	if b.isSquareAttacked(InvalidCoord, White) {
		t.Errorf("expected isSquareAttacked(InvalidCoord, White) == false")
	}
	if b.isSquareAttacked(InvalidCoord, Black) {
		t.Errorf("expected isSquareAttacked(InvalidCoord, Black) == false")
	}
}

func TestGetMove(t *testing.T) {
	tests := []struct {
		b         *Board
		from, to  string
		isCapture bool
	}{
		{New(), "e2", "e4", false},
	}

	for i, test := range tests {
		from, err := CoordFromString(test.from)
		if err != nil {
			t.Errorf("[%d] error making coord: %v", i, from)
		}
		to, err := CoordFromString(test.to)
		if err != nil {
			t.Errorf("[%d] error making coord: %v", i, to)
		}
		move, err := test.b.GetMove(from, to)
		if err != nil {
			t.Errorf("[%d] unexpected error: %s-%s err: %v", i, test.from, test.to, err)
		}
		if move.isCapture != test.isCapture {
			t.Errorf("[%d] (%v).isCapture != %v", i, move, test.isCapture)
		}
	}
}

var moveCounts = []int{20, 400, 8902, 197281, 4865609} //, 119060324}

func TestCountMoves(t *testing.T) {
	/*
		b := New()
		for i, target := range moveCounts {
			t.Logf("%v %v\n", i, target)
			if cnt := b.Perft(i+1, false); cnt != target {
				t.Errorf("[%d] wrong move count: %d != %d", i+1, cnt, target)
				b.Perft(i+1, true)
			}
		}
	*/
	coord := testingCoordFunc(t)
	b := New()
	b.MakeMove(Move{p: White | Pawn, from: coord("f2"), to: coord("f4")})
	b.MakeMove(Move{p: Black | Pawn, from: coord("e7"), to: coord("e5")})
	b.Print()
	b.MakeMove(Move{p: White | Pawn, from: coord("h2"), to: coord("h4")})
	b.Print()
	b.MakeMove(Move{p: Black | Queen, from: coord("d8"), to: coord("h4"), isCapture: true, captured: White | Pawn})
	b.Print()
	fmt.Println(b)
	fmt.Println(b.Perft(1, true))
	m := b.PossibleMoves()
	fmt.Println(len(m), m)
}

func TestPseudoMoves(t *testing.T) {
	b := EmptyBoard()
	for _, p := range []Piece{Pawn, Knight, Bishop, Rook, Queen, King} {
		for _, c := range []Piece{White, Black} {
			for i := 0; i < 64; i++ {
				coord := CoordFromIdx(i)

				// Pawns might not always attack something.
				if p == Pawn {
					if c == White && coord.y == 7 ||
						c == Black && coord.y == 0 {
						continue
					}
				}

				b.set(p|c, coord)
				if b.PseudoMoves(c).countZeros() == 0 {
					t.Errorf("error setting: %v %v", p|c, coord)
				}

				b.set(Empty, coord)
				if b.PseudoMoves(c).countZeros() != 0 {
					t.Errorf("error emptying: %v %v", p|c, coord)
				}
			}
		}
	}
}

func TestCaptureClearsPseudo(t *testing.T) {
	coord := testingCoordFunc(t)
	b := EmptyBoard()
	b.set(Pawn|White, coord("c2"))
	b.set(Pawn|Black, coord("d7"))
	b.MakeMove(Move{p: White | Pawn, from: coord("c2"), to: coord("c4")})
	b.MakeMove(Move{p: Black | Pawn, from: coord("d7"), to: coord("d5")})
	b.MakeMove(Move{p: White | Pawn, from: coord("c4"), to: coord("d5"), isCapture: true, captured: Black | Pawn})
	if l := b.PseudoMoves(White).Attackers(coord("d5")); len(l) != 0 {
		t.Errorf("no white pieces should be attacking: %v", l)
	}
}

func BenchmarkPerft1(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(1, false)
	}
}
func BenchmarkPerft2(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(2, false)
	}
}
func BenchmarkPerft3(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(3, false)
	}
}
func BenchmarkPerft4(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(4, false)
	}
}

func BenchmarkPerft5(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(5, false)
	}
}
