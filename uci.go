package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func uciWriteln(s string) {
	fmt.Println(s)
}

func listOptions() {
	uciWriteln("id name GopherChess")
	uciWriteln("id author Jeremy Faller")
	uciWriteln("uciok")
}

func isReady() {
	uciWriteln("readyok")
}

func setOption(option []string, b *Board) *Board {
	if len(option) == 0 {
		uciWriteln("No such option:")
		return b
	}
	switch option[0] {
	default:
		uciWriteln("No such option")
	}
	return b
}

func uciNewGame() *Board {
	return New()
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

func UCI() error {
	scanner := bufio.NewScanner(os.Stdin)
	b := uciNewGame()
	for scanner.Scan() {
		str := strings.Trim(scanner.Text(), " \t")
		if len(str) == 0 {
			continue
		}
		strs := removeBlanks(strings.Split(str, " "))
		switch strs[0] {
		case "uci":
			listOptions()
		case "isready":
			isReady()
		case "setoption":
			b = setOption(strs[1:], b)
		case "ucinewgame":
			b = uciNewGame()
		default:
			uciWriteln(fmt.Sprintf("Unknown command: %s", str))
		}
	}
	return scanner.Err()
}
