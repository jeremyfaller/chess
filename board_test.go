package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func capture(m Move) Move {
	m.isCapture = true
	return m
}

func setCheck(m Move) Move {
	m.isCheck = true
	return m
}

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
		{ // Capturing can stop a check.
			"capture can stop check",
			"rnb1kbnr/pppp1ppp/8/4p3/5P1q/8/PPPPP1P1/RNBQKBNR w KQkq - 0 3",
			Move{p: White | Rook, to: coord("h4"), from: coord("h1"), isCapture: true},
			true,
			true,
		},
		{ // Moving can stop a check.
			"move can stop check",
			"rnb1kbnr/pp1ppppp/8/q1p5/5B2/3P4/PPP1PPPP/RN1QKBNR w KQkq - 2 3",
			Move{p: White | Bishop, to: coord("d2"), from: coord("f4")},
			true,
			false,
		},
		{ // Can't castle queenside when there's a knight.
			"can't O-O-O with knight",
			"r3k2r/8/8/8/8/8/8/RN2K2R b KQkq - 0 1",
			Move{p: White | King, to: coord("c1"), from: coord("e1")},
			false,
			false,
		},
		{ // Can't capture into check.
			"can't capture into check",
			"r1bqkbnr/1ppppppp/8/p7/1n6/2KP4/PPP1PPPP/RNBQ1BNR w kq - 4 4",
			Move{p: White | King, to: coord("b4"), from: coord("c3"), isCapture: true},
			false,
			true,
		},
		{
			"can castle",
			"r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/1p2P3/2N2Q1p/PPPBBPPP/1R2K2R w Kkq - 2 2",
			Move{p: White | King, from: coord("e1"), to: coord("g1")},
			true,
			false,
		},
	}

	for i, test := range tests {
		t.Logf("[%d] %s", i, test.desc)
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Errorf("[%d] %s error creating board: %v", i, test.desc, err)
		}
		if isLegal := b.isLegalMove(&test.move); isLegal != test.isLegal {
			t.Errorf("[%d] %s isLegalMove = %t, expected %t", i, test.desc, isLegal, test.isLegal)
		}
		if test.move.isCapture != test.isCapture {
			t.Errorf("[%d] %s move.isCapture = %t, expected %t", i, test.desc, test.move.isCapture, test.isCapture)
		}
	}
}

func TestGetMoves(t *testing.T) {
	coord := testingCoordFunc(t)
	move := func(p Piece, from, to string, capPiece Piece) Move {
		if capPiece == Empty {
			return Move{p: p, from: coord(from), to: coord(to)}
		}
		return Move{p: p, from: coord(from), to: coord(to), isCapture: true}
	}
	promo := func(p Piece, from, to string, promo Piece) Move {
		return Move{p: p, from: coord(from), to: coord(to), promotion: promo}
	}
	check := func(m Move) Move {
		m.isCheck = true
		return m
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
				move(White|Pawn, "a2", "a4", Empty),
				move(White|Pawn, "a2", "a3", Empty),
			}),
		},
		{ // Test that a non-opening pawn move is only one space.
			"pawn move",
			"k7/8/8/8/8/P7/8/K7 w - - 0 1",
			coord("a3"),
			toSet([]Move{
				move(White|Pawn, "a3", "a4", Empty),
			}),
		},
		{ // Test moving into your own piece stops the move.
			"same color stops",
			"1k6/8/8/8/8/R7/8/R1K5 w - - 0 1",
			coord("a1"),
			toSet([]Move{
				move(White|Rook, "a1", "a2", Empty),
				check(move(White|Rook, "a1", "b1", Empty)),
			}),
		},
		{
			"king moves",
			"1k6/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			coord("e1"),
			toSet([]Move{
				move(White|King, "e1", "d1", Empty),
				move(White|King, "e1", "d2", Empty),
				move(White|King, "e1", "e2", Empty),
				move(White|King, "e1", "f1", Empty),
				move(White|King, "e1", "f2", Empty),
				move(White|King, "e1", "g1", Empty),
				move(White|King, "e1", "c1", Empty),
			}),
		},
		{
			"promotions",
			"k7/4P3/8/8/8/8/8/K7 w - - 0 1",
			coord("e7"),
			toSet([]Move{
				check(promo(White|Pawn, "e7", "e8", White|Queen)),
				promo(White|Pawn, "e7", "e8", White|Knight),
				promo(White|Pawn, "e7", "e8", White|Bishop),
				check(promo(White|Pawn, "e7", "e8", White|Rook)),
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
		{
			"capture can stop check",
			"rnb1kbnr/pppp1ppp/8/4p3/5P1q/8/PPPPP1P1/RNBQKBNR w KQkq - 0 3",
			coord("h1"),
			toSet([]Move{move(White|Rook, "h1", "h4", Black|Queen)}),
		},
		{
			"bishop can stop check",
			"rnb1kbnr/pp1ppppp/8/q1p5/5B2/3P4/PPP1PPPP/RN1QKBNR w KQkq - 2 3",
			coord("f4"),
			toSet([]Move{move(White|Bishop, "f4", "d2", Empty)}),
		},
	}

	for i, test := range tests {
		t.Logf("[%d] %s %+v\n", i, test.desc, test.c)
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%s] error creating board; %v", test.desc, err)
		}
		moves := toSet(b.GetMoves(nil, test.c))
		if diff := pretty.Compare(test.moves, moves); diff != "" {
			t.Errorf("[%s] moves unequal:\n%s", test.desc, diff)
		}
	}
}

func TestFENChecks(t *testing.T) {
	tests := []struct {
		desc    string
		fen     string
		isCheck bool
	}{
		{
			"no check",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			false,
		},
		{
			"white check",
			"1k6/8/8/8/8/4r3/8/R3K2R w - - 0 1",
			true,
		},
		{
			"black check",
			"r3k2r/8/4R3/8/8/8/8/1K6 b - - 0 1",
			true,
		},
		{
			"no pawn check",
			"k7/P7/8/8/8/8/8/K7 b - - 0 1",
			false,
		},
		{
			"pawn check",
			"k7/1P6/8/8/8/8/8/K7 b - - 0 1",
			true,
		},
	}

	for _, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%s] unexpected error: %v", test.desc, err)
		}
		if v := b.IsCheck(); v != test.isCheck {
			t.Errorf("[%s] IsCheck() = %v, expected = %v", test.desc, v, test.isCheck)
		}
	}
}

func TestWouldKingBeInCheck(t *testing.T) {
	coord := testingCoordFunc(t)
	move := func(p Piece, from, to string, isCap bool) Move {
		return Move{p: p, from: coord(from), to: coord(to), isCapture: isCap}
	}

	tests := []struct {
		desc  string
		fen   string
		move  Move
		would bool
	}{
		{
			"rook captures queen",
			"rnb1kbnr/pppp1ppp/8/4p3/5P1q/8/PPPPP1P1/RNBQKBNR w KQkq - 0 3",
			move(White|Rook, "h1", "h4", true),
			false,
		},
		{
			"bishop can stop check",
			"rnb1kbnr/pp1ppppp/8/q1p5/5B2/3P4/PPP1PPPP/RN1QKBNR w KQkq - 2 3",
			move(White|Bishop, "f4", "d2", false),
			false,
		},
	}

	for i, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Errorf("[%d] error converting fen(%q), %v", i, test.fen, err)
		}
		if v := b.wouldKingBeInCheck(&test.move); v != test.would {
			t.Errorf("[%d] wouldKingBeInCheck(%v) = %v, expected %v", i, test.move, v, test.would)
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
		t.Logf("%s", test.desc)
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

func TestPerft(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping perft tests")
	}

	/*
		coord := testingCoordFunc(t)
		move := func(p Piece, from, to string) Move {
			return Move{p: p, from: coord(from), to: coord(to)}
		}
		check := func(m Move) Move { m.isCheck = true; return m }
	*/
	tests := []struct {
		desc   string
		fen    string
		depth  int
		target uint64
		moves  []Move
	}{
		{"perft(1)", StartingFEN, 1, 20, nil},
		{"perft(2)", StartingFEN, 2, 400, nil},
		{"perft(3)", StartingFEN, 3, 8902, nil},
		{"perft(4)", StartingFEN, 4, 197281, nil},
		{"perft(5)", StartingFEN, 5, 4865609, nil},
		{"perft(6)", StartingFEN, 6, 119060324, nil},
		/*
			{"perft(7)", StartingFEN, 7, 3195901860, nil},
			{"perft(8)", StartingFEN, 8, 84998978956, nil},
		*/
		{"kiwipete(1)", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1, 48, nil},
		{"kiwipete(2)", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2, 2039, nil},
		{"kiwipete(3)", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3, 97862, nil},
		{"kiwipete(4)", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4, 4085603, nil},
		{"kiwipete(5)", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5, 193690690, nil},
		/*
			{"kiwipete(6)", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 6, 8031647685, nil},
		*/
		{"position 5", "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", 5, 89941194, nil},
	}

	for i, test := range tests {
		test := test // rebind
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Logf("[%d] %v should equal %d", i, test.desc, test.target)
			b, err := FromFEN(test.fen)
			if err != nil {
				t.Errorf("[%d] %s FromFEN(%q) = %v", i, test.desc, test.fen, err)
			}
			if len(test.moves) != 0 {
				t.Logf("[%d] %s executing moves %v", i, test.desc, test.moves)
				for _, move := range test.moves {
					b.MakeMove(move)
				}
			}
			if cnt := b.Perft(test.depth - len(test.moves)); cnt != test.target {
				t.Errorf("[%d] %s perft(%v) = %d, expected %d", i, test.desc, b.FENString(), cnt, test.target)
			}
		})
	}
}

func TestPossibleMoves(t *testing.T) {
	coord := testingCoordFunc(t)
	move := func(p Piece, from, to string) Move {
		return Move{p: p, from: coord(from), to: coord(to)}
	}
	tests := []struct {
		desc  string
		fen   string
		moves []Move
	}{
		{
			"black block check",
			"rnbqkbnr/ppp1pppp/8/3p4/Q7/2P5/PP1PPPPP/RNB1KBNR b KQkq - 1 2",
			[]Move{move(Black|Pawn, "c7", "c6"), move(Black|Pawn, "b7", "b5"), move(Black|Knight, "b8", "c6"),
				move(Black|Knight, "b8", "d7"), move(Black|Bishop, "c8", "d7"), move(Black|Queen, "d8", "d7")},
		},
		{
			"white block check",
			"rnb1kbnr/pp1ppppp/8/q1p5/5B2/3P4/PPP1PPPP/RN1QKBNR w KQkq - 2 3",
			[]Move{move(White|Knight, "b1", "d2"), move(White|Knight, "b1", "c3"), move(White|Pawn, "b2", "b4"),
				move(White|Pawn, "c2", "c3"), move(White|Queen, "d1", "d2"), move(White|Bishop, "f4", "d2")},
		},
	}

	for i, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("[%d] %s bad fen %v", i, test.desc, err)
		}
		moves := b.PossibleMoves(nil)
		if len(moves) != len(test.moves) {
			t.Logf("%v != %v", moves, test.moves)
			t.Errorf("[%d] %s, len(moves) = %d, expected %d", i, test.desc, len(moves), len(test.moves))
		}
		for _, move := range test.moves {
			found := false
			for _, try := range moves {
				if try == move {
					found = true
					break
				}
			}
			if found == false {
				t.Errorf("[%d] %s move: %v, not found", i, test.desc, move)
			}
		}
	}
}

func TestQueenMoves(t *testing.T) {
	b, err := FromFEN("k7/8/8/3Q4/8/8/8/K7 w - - 1 1")
	p := Piece(White | Queen)
	coord := testingCoordFunc(t)
	moves := []Move{
		Move{p: p, from: coord("d5"), to: coord("h1"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("a2"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("g2"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("b3"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("f3"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("c4"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("e4"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("c6"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("e6"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("b7"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("f7"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("g8"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("d1"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("d2"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("d3"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("d4"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("a5"), isCheck: true},
		Move{p: p, from: coord("d5"), to: coord("b5"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("c5"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("e5"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("f5"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("g5"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("h5"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("d6"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("d7"), isCheck: false},
		Move{p: p, from: coord("d5"), to: coord("d8"), isCheck: true},
	}
	if err != nil {
		t.Errorf("error making board: %v", err)
	}
	if res := b.GetMoves(nil, coord("d5")); !reflect.DeepEqual(moves, res) {
		t.Errorf("GetMoves(nil, d5) = %v, expected = %v", res, moves)
	}
}

func TestOccupancy(t *testing.T) {
	tests := []struct {
		fen string
		occ Bit
	}{
		{"k7/8/8/8/8/r7/R7/K7 b - - 0 1", 0x0100000000010101},
		{"7k/8/8/8/8/r7/R7/K7 b - - 0 1", 0x8000000000010101},
	}

	for _, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Errorf("error making board: %v", err)
		}
		if occ := b.state.wOcc | b.state.bOcc; occ != test.occ {
			t.Errorf("expected occupancy = %v, got %v", test.occ, occ)
		}
	}
}

func TestSimpleScore(t *testing.T) {
	tests := []struct {
		p     Piece
		score Score
	}{
		{Piece(White | Pawn), 100},
		{Piece(Black | Pawn), -100},
		{Piece(White | Bishop), 300},
		{Piece(Black | Bishop), -300},
		{Piece(White | Knight), 300},
		{Piece(Black | Knight), -300},
		{Piece(White | Rook), 500},
		{Piece(Black | Rook), -500},
		{Piece(White | Queen), 900},
		{Piece(Black | Queen), -900},
	}

	for i, test := range tests {
		b := EmptyBoard()
		b.set(test.p, CoordFromIdx(0))
		if s := b.state.score; s != test.score {
			t.Errorf("[%d] b.Score() = %d, exepcted = %d", i, s, test.score)
		}
	}
}

func TestHashDifference(t *testing.T) {
	tests := []struct {
		f1, f2 string
	}{
		{"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1", "rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR b KQkq - 0 1"},
		{"rnbqkb1r/pppppppp/5n2/8/8/7P/PPPPPPP1/RNBQKBNR w KQkq - 1 2", "rnbqkb1r/pppppppp/5n2/8/7P/8/PPPPPPP1/RNBQKBNR w KQkq - 1 2"},
	}

	for i, test := range tests {
		if test.f1 == test.f2 {
			t.Errorf("[%d] fen strings equal", i)
		}
		b1, err := FromFEN(test.f1)
		if err != nil {
			t.Fatalf("[%d] FromFEN(%q) = %v", i, test.f1, err)
		}
		b2, err := FromFEN(test.f2)
		if err != nil {
			t.Fatalf("[%d] FromFEN(%q) = %v", i, test.f2, err)
		}
		if b1.state.hash == b2.state.hash {
			t.Errorf("[%d] hashes are equal", i)
		}
	}
}

func TestHashSame(t *testing.T) {
	b1 := New()
	b2, err := FromFEN(StartingFEN)
	if err != nil {
		t.Fatalf("error making board from starting FEN: %v", err)
	}
	if h1, h2 := b1.ZHash(), b2.ZHash(); h1 != h2 {
		t.Errorf("hash %x != %x", h1, h2)
	}
}

func TestBoardWithMoves(t *testing.T) {
	tests := []struct {
		fen   string
		moves []string
		isErr bool
		res   string
	}{
		{"k7/8/8/8/8/8/8/K7 w - - 0 1", []string{"a1a2"}, false, "k7/8/8/8/8/8/K7/8 b - - 1 1"},
		{"k7/8/8/8/8/8/8/K7 b - - 1 1", []string{"a8a7"}, false, "8/k7/8/8/8/8/8/K7 w - - 2 2"},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", []string{"a1a2"}, false, "r3k2r/8/8/8/8/8/R7/4K2R b Kkq - 1 1"},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", []string{"e1g1"}, false, "r3k2r/8/8/8/8/8/8/R4RK1 b kq - 1 1"},
		{"k7/1Q6/8/8/8/8/8/K7 w - - 0 1", []string{"b7b8q"}, false, "kQ6/8/8/8/8/8/8/K7 b - - 1 1"},
	}

	for i, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Errorf("[%d] error creating board: %v", i, err)
		}
		if err := b.ApplyMoves(test.moves); err != nil && test.isErr == false {
			t.Errorf("[%d] b.ApplyMoves(%v) = %v, expected nil", i, test.moves, err)
		}
		if fen := b.FENString(); fen != test.res {
			t.Errorf("[%d] resultant FEN = %q, expected %q", i, fen, test.res)
		}
	}
}

func TestHasEnoughMaterialForMate(t *testing.T) {
	tests := []struct {
		desc   string
		pieces []Piece
		isMate bool
	}{
		{"only kings", []Piece{}, false},
		{"with knight", []Piece{White | Knight}, false},
		{"with 2xknight", []Piece{White | Knight, White | Knight}, false},
		{"with 2xknight and extra", []Piece{White | Knight, White | Knight, Black | Knight}, true},
		{"pawn", []Piece{White | Pawn}, true},
		{"bishops", []Piece{Black | Bishop, Black | Bishop}, true},
		{"rook", []Piece{White | Rook}, true},
		{"queen", []Piece{Black | Queen}, true},
	}

	for _, test := range tests {
		b := EmptyBoard()
		b.set(White|King, CoordFromIdx(0))
		b.set(Black|King, CoordFromIdx(1))
		for i := range test.pieces {
			b.set(test.pieces[i], CoordFromIdx(i+2))
		}
		if v := b.hasEnoughMaterialForMate(); v != test.isMate {
			t.Errorf("[%v] hasEnoughMaterialForMate = %t, expected %t", test.desc, v, test.isMate)
		}
	}
}

func TestIsDrawn(t *testing.T) {
	tests := []struct {
		fen   string
		moves []string
		drawn bool
	}{
		{StartingFEN, []string{"b1c3", "b8c6", "c3b1", "c6b8", "b1c3", "b8c6", "c3b1", "c6b8"}, true},
		{strings.Replace(StartingFEN, "0 1", "49 1", 1), []string{"b1c3"}, true},
		{strings.Replace(StartingFEN, "0 1", "49 1", 1), []string{"a2a4"}, false},
		{"K7/8/8/8/8/8/8/k7 - - - 0 1", []string{}, true},
		{"KR6/8/8/8/8/8/8/k7 - - - 0 1", []string{}, false},
		{"KQ6/8/8/8/8/8/8/k7 - - - 0 1", []string{}, false},
		{"KBB5/8/8/8/8/8/8/k7 - - - 0 1", []string{}, false},
		{"KNB5/8/8/8/8/8/8/k7 - - - 0 1", []string{}, false},
		{"K7/P7/8/8/8/8/8/k7 - - - 0 1", []string{}, false},
		{"KNN5/8/8/8/8/8/8/k7 - - - 0 1", []string{}, true},
		{"KNN5/8/8/8/8/8/n7/k7 - - - 0 1", []string{}, false},
	}

	for i, test := range tests {
		b, err := FromFEN(test.fen)
		if err != nil {
			t.Fatalf("error parsing board: %q – %v", test.fen, err)
		}
		for j, m := range test.moves {
			if b.isDrawn() {
				t.Errorf("[%d], isDrawn() == true, expected false", j)
			}
			if err := b.ApplyStringMove(m); err != nil {
				t.Errorf("error applying move: %v %v", m, err)
			}
		}
		if b.isDrawn() != test.drawn {
			t.Errorf("[%d] isDrawn() == %t, expected %t", i, b.isDrawn(), test.drawn)
		}
	}
}

func BenchmarkPerft1(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(1)
	}
}
func BenchmarkPerft2(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(2)
	}
}
func BenchmarkPerft3(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(3)
	}
}
func BenchmarkPerft4(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(4)
	}
}

func BenchmarkPerft5(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(5)
	}
}

func BenchmarkPerft6(b *testing.B) {
	board := New()
	for n := 0; n < b.N; n++ {
		board.Perft(6)
	}
}
