package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	welcome = `GopherChess by Jeremy Faller (jeremy.faller@gmail.com)`
	ws      = " \t"
)

func (u *UCI) defaultContext() context.Context {
	u.ctx, u.cancel = context.WithCancel(context.Background())
	return u.ctx
}

type UCI struct {
	b      *Board
	ctx    context.Context
	cancel context.CancelFunc
}

func trim(s string) string {
	return strings.Trim(s, ws)
}

func removeBlanks(strs []string) []string {
	ret := []string{}
	for _, s := range strs {
		if len(s) != 0 {
			ret = append(ret, s)
		}
	}
	return ret
}

func (u *UCI) Writeln(s string) {
	fmt.Println(s)
}

func (u *UCI) listOptions() {
	u.Writeln("id name GopherChess")
	u.Writeln("id author Jeremy Faller (jeremy.faller@gmail.com)")
	u.Writeln("uciok")
}

func (u *UCI) isReady() {
	u.Writeln("readyok")
}

func (u *UCI) setOption(option []string) {
	if len(option) == 0 {
		u.Writeln("No such option:")
	}
	switch option[0] {
	default:
		u.Writeln("No such option")
	}
}

// position sets the position for the chess engine.
func (u *UCI) position(cmd string) error {
	// Split up the string.
	var moves []string
	fen, moveStr, found := strings.Cut(cmd, "moves")
	if found {
		moves = strings.Split(trim(moveStr), ws)
	}

	// Get the starting position.
	fen = trim(fen)
	if fen == "startpos" {
		fen = StartingFEN
	}
	if b, err := FromFEN(fen); err != nil {
		return err
	} else {
		u.b = b
	}

	// Apply the moves.
	return u.b.ApplyMoves(moves)
}

// newgame creates a new game.
func (u *UCI) newGame() error {
	return u.position("startpos")
}

func (u *UCI) goCmd(cmd string) error {
	kind, opts, _ := strings.Cut(cmd, " ")
	switch kind {
	case "perft":
		res := strings.SplitN(trim(opts), ws, 2)
		cnt, err := strconv.Atoi(res[0])
		if err == nil {
			fmt.Printf("\n\nTotal nodes: %d\n\n", u.b.Perft(cnt))
		} else {
			return fmt.Errorf("no perft count specified")
		}
	}

	return nil
}

func (u *UCI) stopCmd() {
	if u.cancel != nil {
		u.cancel()
	}
}

func (u *UCI) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	if err := u.newGame(); err != nil {
		panic(err)
	}
	u.Writeln(welcome)
	for scanner.Scan() {
		var err error
		str := trim(scanner.Text())
		if len(str) == 0 { // Skip blanks.
			continue
		}

		strs := removeBlanks(strings.Split(str, " "))
		cmdStripped := trim(strings.TrimPrefix(str, strs[0]))
		switch strs[0] {
		case "uci":
			u.listOptions()
		case "isready":
			u.isReady()
		case "setoption":
			u.setOption(strs[1:])
		case "ucinewgame":
			err = u.newGame()
		case "position":
			err = u.position(cmdStripped)
		case "go":
			err = u.goCmd(cmdStripped)
		case "stop":
			u.stopCmd()
		default:
			err = fmt.Errorf("Unknown command: %s", str)
		}

		if err != nil {
			u.Writeln(fmt.Sprintf("%v", err))
		}
	}
	return scanner.Err()
}
