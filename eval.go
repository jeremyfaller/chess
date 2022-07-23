package main

import (
	"fmt"
	"time"
)

// line holds the moves and the score for the current line.
type line struct {
	s Score
	m []Move
}

type lineQ []line

func (q lineQ) Len() int {
	return len(q)
}

func (q lineQ) Less(i, j int) bool {
	return q[i].s < q[j].s
}

func (q lineQ) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *lineQ) Push(l any) {
	*q = append(*q, l.(line))
}

func (q *lineQ) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

type Eval struct {
	b         *Board
	positions int
	l         line
	lq        *lineQ

	// benchmark evaluations
	startTime, endTime time.Time
}

// Creates a new Eval.
func NewEval(b *Board) Eval {
	return Eval{
		b:  b,
		lq: &lineQ{},
	}
}

// depth returns how deep we should run our evaluation.
func (e *Eval) depth() int {
	return 3
}

// calc evaluates the current position, and returns a score.
func (e *Eval) calc(player Piece) Score {
	return e.b.CurrentPlayerScore()
}

// Start begins an evaluation.
func (e *Eval) Start() {
	origDepth := e.depth()
	movesToCheck := make([][]Move, origDepth+1)
	player := e.b.state.turn
	e.b.Print()

	var search func(int, Score, Score) Score
	search = func(d int, alpha, beta Score) Score {
		// Stats.
		e.positions += 1

		// Get a link to our local slice.
		moves := movesToCheck[origDepth-d][:0]
		moves = e.b.PossibleMoves(moves)
		// If no moves, we could be in stalemate or checkmate.
		if len(moves) == 0 {
			if e.b.IsKingInCheck(e.b.state.turn) {
				return checkmate
			}
			return stalemate
		}

		// If we're done, just calculate the score.
		// TODO(jfaller): Might want to do this after completing all
		// checks or captures.
		if d == 0 {
			return e.calc(player)
		}

		// Alpha-beta prune the search tree.
		if d == origDepth {
			fmt.Println(moves)
		}
		for _, move := range moves {
			e.b.MakeMove(move)
			evaluation := -search(d-1, -beta, -alpha)
			e.b.UnmakeMove()

			// Prune early.
			// TODO(jfaller): need to check that the engine will find mulitple
			// checkmates from a given position.
			if evaluation >= beta {
				return beta
			}
			if evaluation > alpha {
				alpha = evaluation
			}
		}
		return alpha
	}

	e.startTime = time.Now()
	score := search(origDepth, minScore, maxScore)
	e.endTime = time.Now()
	fmt.Println(score)
}
