package main

import (
	"fmt"
	"time"
)

type Eval struct {
	b         *Board
	positions int
	depth     int
	score     Score

	// benchmark evaluations
	startTime, endTime time.Time
}

// Creates a new Eval.
func NewEval(b *Board) Eval {
	e := Eval{
		b:     b,
		depth: 3,
	}
	return e
}

// sortMoves sorts the possible moves, trying to find good ones first.
func (e *Eval) sortMoves(moves []Move) {
	// Move checks and captures to the head of the queue.
	pIdx := 0
	for i := range moves {
		if moves[i].isCheck || moves[i].isCapture || moves[i].IsPromotion() {
			moves[pIdx], moves[i] = moves[i], moves[pIdx]
			pIdx += 1
		}
	}
}

// calc evaluates the current position, and returns a score.
func (e *Eval) calc(player Piece) Score {
	return e.b.CurrentPlayerScore()
}

// Duration returns the length of time the evaluation has run.
func (e *Eval) Duration() time.Duration {
	return e.endTime.Sub(e.startTime)
}

// TimeString returns the length of time it took to do the evaluation.
func (e *Eval) TimeString() string {
	d := e.Duration()
	if d.Hours() > 1 {
		return fmt.Sprintf("%gh", d.Hours())
	}
	if d.Minutes() > 1 {
		return fmt.Sprintf("%gm", d.Minutes())
	}
	if d.Seconds() > 1 {
		return fmt.Sprintf("%gs", d.Seconds())
	}
	if d.Milliseconds() > 1 {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d.Microseconds() > 1 {
		return fmt.Sprintf("%dus", d.Microseconds())
	}
	return fmt.Sprintf("%dns", d.Nanoseconds())
}

// Start begins an evaluation.
func (e *Eval) Start() {
	origDepth := e.depth
	movesToCheck := make([][]Move, origDepth+1)
	player := e.b.state.turn

	var search func(int, Score, Score) Score
	search = func(d int, alpha, beta Score) Score {
		// Stats.
		e.positions += 1

		// Get a link to our local slice.
		moves := movesToCheck[origDepth-d][:0]
		moves = e.b.PossibleMoves(moves)
		e.sortMoves(moves)

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
	e.score = search(origDepth, minScore, maxScore)
	e.endTime = time.Now()
}
