package main

import (
	"sync"
	"sync/atomic"
)

type TTType uint8

const (
	TTExact TTType = iota
	TTUpper
	TTLower
)

type TTStats struct {
	Entries uint64
	Inserts uint64
	Lookups uint64
	Hits    uint64
	Misses  uint64
}

// ttEntry is an entry in a TranspositionTable.
type ttEntry struct {
	hash  Hash
	move  Move
	score Score
	depth Depth
	t     TTType
}

// TranspositionTable holds transpositions.
type TranspositionTable struct {
	m    sync.RWMutex
	vals []ttEntry

	// stats
	lookups atomic.Uint64
	misses  atomic.Uint64
	inserts atomic.Uint64
}

// NewTranspositionTable creates a new TranspositionTable of a given size.
func NewTranspositionTable(sizeMB int) *TranspositionTable {
	tt := &TranspositionTable{}
	tt.Resize(sizeMB)
	return tt
}

// Resize resizes the TranspositionTable.
func (tt *TranspositionTable) Resize(sizeMB int) {
	if sizeMB < 0 {
		panic("invalid transposition table size")
	}

	tt.m.Lock()
	defer tt.m.Unlock()

	// Resize the table.
	entries := (sizeMB * 1024 * 1024) / getSize(ttEntry{})
	tt.vals = make([]ttEntry, entries)
	tt.clearStats()
}

// Clear removes all entries from a TranspositionTable.
func (tt *TranspositionTable) Clear() {
	tt.m.Lock()
	defer tt.m.Unlock()

	clear(tt.vals)
	tt.clearStats()
}

// Size returns the number of entries in the TranspositionTable.
func (tt *TranspositionTable) Size() int {
	return len(tt.vals)
}

// Entries counts the number of non-zero entries.
func (tt *TranspositionTable) Entries() (count int) {
	tt.m.RLock()
	defer tt.m.RUnlock()

	for i := range tt.vals {
		if tt.vals[i].hash != 0 {
			count += 1
		}
	}
	return count
}

// clearStats clears the stats.
func (tt *TranspositionTable) clearStats() {
	tt.lookups.Store(0)
	tt.misses.Store(0)
	tt.inserts.Store(0)
}

// Stats returns the stats structure.
func (tt *TranspositionTable) Stats() TTStats {
	l, m, i := tt.lookups.Load(), tt.misses.Load(), tt.inserts.Load()
	return TTStats{
		Entries: uint64(tt.Entries()),
		Lookups: l,
		Hits:    l - m,
		Misses:  m,
		Inserts: i,
	}
}

// index returns the index of the given Hash.
func (tt *TranspositionTable) index(hash Hash) int {
	return int(hash % Hash(len(tt.vals)))
}

// Lookup tries to find an entry in the TranspositionTable.
func (tt *TranspositionTable) Lookup(hash Hash, depth, plyRemain Depth, alpha, beta Score) (Score, Move, bool) {
	tt.m.RLock()
	defer tt.m.RUnlock()

	tt.lookups.Add(1)
	entry := tt.vals[tt.index(hash)]
	if entry.hash == hash && entry.depth >= plyRemain {
		score := correctScore(entry.score, depth)
		if entry.t == TTExact {
			return score, entry.move, true
		}
		if entry.t == TTUpper && score <= alpha {
			return score, entry.move, true
		}
		if entry.t == TTLower && score >= beta {
			return score, entry.move, true
		}
	}
	tt.misses.Add(1)
	return 0, Move{}, false
}

// Insert puts an entry into the transposition table.
func (tt *TranspositionTable) Insert(hash Hash, move Move, score Score, plySearched, plyRemain Depth, evalType TTType) {
	tt.m.Lock()
	defer tt.m.Unlock()

	tt.inserts.Add(1)
	tt.vals[tt.index(hash)] = ttEntry{
		hash:  hash,
		move:  move,
		score: correctScore(score, plySearched),
		depth: plyRemain,
		t:     evalType,
	}
}

// correctScore
func correctScore(score Score, numPly Depth) Score {
	if IsMateScore(score) {
		sign := Score(1)
		if score < 0 {
			sign = -1
		}
		return (score*sign - Score(numPly)) * sign
	}
	return score
}
