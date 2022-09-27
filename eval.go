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

type Line struct {
	c     Piece
	Score Score
	Moves []Move
}

func newLine(d int, c Piece) Line {
	return Line{
		c:     c,
		Score: minScore,
		Moves: make([]Move, d),
	}
}

func (l *Line) add(m Move, s Score, d int) {
	if m.p.Color() != l.c {
		s *= -1
	}

	// Ignore worse Lines.
	if s < l.Score && d != 0 {
		return
	}

	l.Score = s
	if len(l.Moves) < d+1 {
		newMoves := make([]Move, d+1)
		copy(newMoves, l.Moves)
		l.Moves = newMoves
	}
	l.Moves[d] = m
}

// Creates a new Eval.
func NewEval(b *Board, depth int) Eval {
	e := Eval{
		b:     b,
		depth: depth,
	}
	return e
}

// sortMoves sorts the possible moves, trying to find good ones first.
//
// We also return the number of moves that are special, ie checks, promotions,
// and captures.
///
// Moves will be sorted in the following order:
//   [0..X] Checks
//   [X..Y] Promotions
//   [Y..Z] Captures
//   [Z..N] Rest
func (e *Eval) sortMoves(moves []Move) int {
	// Move likely good moves to the head.
	idx := 0
	for i := range moves {
		if moves[i].isCheck || moves[i].isCapture || moves[i].IsPromotion() {
			moves[idx], moves[i] = moves[i], moves[idx]
			idx += 1
		}
	}

	// Move checks to the begining.
	cIdx := 0
	for i := 0; i < idx; i++ {
		if moves[i].isCheck {
			moves[cIdx], moves[i] = moves[i], moves[cIdx]
			cIdx += 1
		}
	}

	// Move promotions to the next place.
	pIdx := 0
	for i := cIdx; i < idx; i++ {
		if moves[i].IsPromotion() {
			moves[pIdx], moves[i] = moves[i], moves[pIdx]
			pIdx += 1
		}
	}
	return idx
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
	line := newLine(origDepth, player)

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

			// Add the move to the line.
			line.add(move, evaluation, origDepth-d)
			if d == origDepth {
				line = newLine(origDepth, player)
			}

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
