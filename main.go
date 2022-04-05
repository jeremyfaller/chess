package main

import "fmt"

func main() {
	b := New()
	b.Print()
	fmt.Println(b.String())
	fmt.Printf("%+v\n", b.PossibleMoves())
}
