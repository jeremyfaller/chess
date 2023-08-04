package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var (
	welcome       = `GopherChess by Jeremy Faller (jeremy.faller@gmail.com)`
	ws            = " \t"
	optionErr     = "No such option"
	unknownCmdErr = "Unknown command"
)

type UCI struct {
	e *Eval
	b *Board
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
	n := runtime.GOMAXPROCS(0)
	u.Writeln("id name GopherChess")
	u.Writeln("id author Jeremy Faller (jeremy.faller@gmail.com)")
	u.Writeln("")
	u.Writeln(fmt.Sprintf("option name Threads type spin default %d min 1 max %d", n, n))
	u.Writeln("uciok")
}

func (u *UCI) isReady() {
	u.Writeln("readyok")
}

func (u *UCI) printError(str string, cmds []string) {
	u.Writeln(fmt.Sprintf("%s: %q", str, strings.Join(cmds, " ")))
}

func (u *UCI) setOption(tokens []string) {
	if len(tokens) != 3 && tokens[1] != "value" {
		u.printError(optionErr, tokens)
		return
	}

	switch tokens[0] {
	case "Thread":
		if v, err := strconv.Atoi(tokens[2]); err != nil {
			u.printError(optionErr, tokens)
		} else {
			runtime.GOMAXPROCS(v)
		}

	default:
		u.printError(optionErr, tokens)
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
	u.e.Stop()
}

func (u *UCI) debug(opts []string) {
	if len(opts) != 1 {
		u.printError(unknownCmdErr, opts)
		return
	}
	switch opts[0] {
	case "on":
		u.e.SetDebug(true)
	case "off":
		u.e.SetDebug(false)
	default:
		u.printError(unknownCmdErr, opts)
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
		case "debug":
			u.debug(strs[1:])
		case "isready":
			u.isReady()
		case "go":
			err = u.goCmd(cmdStripped)
		case "position":
			err = u.position(cmdStripped)
		case "setoption":
			u.setOption(strs[2:])
		case "stop":
			u.stopCmd()
		case "uci":
			u.listOptions()
		case "ucinewgame":
			err = u.newGame()
		default:
			u.printError(unknownCmdErr, strs)
		}

		if err != nil {
			u.Writeln(fmt.Sprintf("%v", err))
		}
	}
	return scanner.Err()
}
