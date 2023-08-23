package main

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"
)

type Depth uint8

const (
	maxScore  = 29999
	minScore  = -maxScore
	stalemate = 0
	checkmate = 10000
)

func IsMateScore(s Score) bool {
	return s+50 > checkmate || s-50 < -checkmate
}

// GameResult signifies what's happening in the game.
type GameResult int

const (
	InProgress GameResult = iota
	Draw
	WhiteIsMated
	BlackIsMated
)

type doneChan chan struct{}

type Eval struct {
	positions int
	depth     Depth
	score     Score
	rand      *rand.Rand

	// benchmark evaluations
	totalTime time.Duration

	// Functions for starting/stopping evaluation.
	m        sync.Mutex
	duration time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool

	tt *TranspositionTable

	// Options
	useBook bool
	debug   bool
}

// Creates a new Eval.
func NewEval(depth Depth) Eval {
	return Eval{
		depth:   depth,
		useBook: true,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		tt:      NewTranspositionTable(20),
	}
}

// SetDebug sets the debug state.
func (e *Eval) SetDebug(v bool) {
	e.debug = v
}

// SetBook determines if the evaluation engine will use the book.
func (e *Eval) SetBook(v bool) {
	e.useBook = v
}

// SetTranspositionTableSize sets the size (in MB) of the TranspositionTable.
func (e *Eval) SetTranspositionTableSize(sizeMB int) {
	e.tt.Resize(sizeMB)
}

// SetDuration stops the current evaluation, and
func (e *Eval) SetDuration(d time.Duration) *Eval {
	if e.running {
		panic("can't SetDuration on a running Eval")
	}
	e.duration = d
	return e
}

// IsRunning returns true if the eval engine is running.
func (e *Eval) IsRunning() bool {
	return e.running
}

// reportMove reports the move.
func (e *Eval) reportMove(m Move) {
	fmt.Println("bestmove", m)
}

// sortMoves sorts the possible moves, trying to find good ones first.
//
// We also return the number of moves that are special, ie checks, promotions,
// and captures.
// /
// Moves will be sorted in the following order:
//
//	[0..X] Checks
//	[X..Y] Promotions
//	[Y..Z] Captures
//	[Z..N] Rest
func (e *Eval) sortMoves(moves []Move, b *Board) int {
	// Find likely good moves.
	idx := 0
	for i, m := range moves {
		if m.isCheck || m.IsPromotion() || m.isCapture {
			moves[idx], moves[i] = moves[i], moves[idx]
			idx += 1
		}
	}

	// Move checks to the begining.
	pos := 0
	for i := pos; i < idx; i++ {
		if moves[i].isCheck {
			moves[pos], moves[i] = moves[i], moves[pos]
			pos += 1
		}
	}

	// Move promotions to the next place.
	for i := pos; i < idx; i++ {
		if moves[i].IsPromotion() {
			moves[pos], moves[i] = moves[i], moves[pos]
			pos += 1
		}
	}

	return idx
}

// calc evaluates the current position, and returns a score.
func (e *Eval) calc(b *Board) Score {
	return b.CurrentPlayerMaterial()
}

// Duration returns the length of time the evaluation has run.
func (e *Eval) Duration() time.Duration {
	return e.totalTime
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

// Stop stops an evaluation.
func (e *Eval) Stop() {
	e.m.Lock()
	defer e.m.Unlock()

	if e.cancel != nil {
		e.cancel()
	}
	e.Wait()
}

// setup creates the context for an evaluation.
func (e *Eval) setup() {
	if e.duration != 0 {
		e.ctx, e.cancel = context.WithTimeout(context.Background(), e.duration)
	} else {
		e.ctx, e.cancel = context.WithCancel(context.Background())
	}
}

// Wait delays until an evaluation is done.
func (e *Eval) Wait() {
	for e.running {
		time.Sleep(time.Millisecond)
	}
}

// Start begins an evaluation.
func (e *Eval) Start(b *Board) {
	// Stop and previously running evaluation.
	e.Stop()

	// And setup for a new one.
	e.m.Lock()
	defer e.m.Unlock()
	e.setup()

	movesToCheck := make([][]Move, e.depth+1)

	shouldCancel := func() bool {
		select {
		case <-e.ctx.Done():
			return true
		default:
			return false
		}
	}

	if e.useBook {
		if move, found := getBook(b, e.rand); found {
			e.reportMove(move)
			return
		}
	}

	line := []Move{}

	var search func(Depth, Depth, Score, Score) Score
	search = func(d, targetD Depth, alpha, beta Score) Score {
		// If we've already seen this position, we don't need to keep searching.
		if ttVal, _, found := e.tt.Lookup(b.ZHash(), d, targetD-d, alpha, beta); found {
			return ttVal
		}
		var bestMove Move
		evalBound := TTUpper

		// Stats.
		e.positions += 1

		// Get a link to our local slice.
		moves := movesToCheck[d][:0]
		moves = b.PossibleMoves(moves)
		e.sortMoves(moves, b)

		// If no moves, we could be in stalemate or checkmate.
		if len(moves) == 0 {
			if b.IsCheck() {
				return -(checkmate - Score(d))
			}
			return stalemate
		}

		// If we're done, just calculate the score.
		// TODO(jfaller): Might want to do this after completing all
		// checks or captures.
		if d == targetD {
			return e.calc(b)
		}

		// Alpha-beta prune the search tree.
		for _, move := range moves {
			if shouldCancel() {
				break
			}

			b.MakeMove(move)
			line = append(line, move)
			evaluation := -search(d+1, targetD, -beta, -alpha)
			line = slices.Delete(line, len(line)-1, len(line))
			b.UnmakeMove()

			// Prune early.
			if evaluation >= beta {
				e.tt.Insert(b.ZHash(), move, beta, d, targetD-d, TTLower)
				return beta
			}
			if evaluation > alpha {
				bestMove = move
				alpha = evaluation
				evalBound = TTExact
			}
		}
		if !bestMove.IsNull() {
			e.tt.Insert(b.ZHash(), bestMove, alpha, d, targetD-d, evalBound)
			if d == 0 {
				e.reportMove(bestMove)
			}
		}
		return alpha
	}

	e.running = true
	go func() {
		startTime := time.Now()
		e.score = search(0, e.depth, minScore, maxScore)
		e.totalTime += time.Now().Sub(startTime)
		e.running = false
	}()
}
