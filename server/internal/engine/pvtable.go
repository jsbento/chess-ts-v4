package engine

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
)

func (e *Engine) StorePvMove(move int) {
	idx := e.Board.PosKey % uint64(e.Board.HashTable.NumEntries)

	e.Board.HashTable.HashEntries[idx].PosKey = e.Board.PosKey
	e.Board.HashTable.HashEntries[idx].Move = move
}

func (e *Engine) ProbePvTable() int {
	idx := e.Board.PosKey % uint64(e.Board.HashTable.NumEntries)

	if e.Board.HashTable.HashEntries[idx].PosKey == e.Board.PosKey {
		return e.Board.HashTable.HashEntries[idx].Move
	}

	return c.NOMOVE
}

func (e *Engine) GetPvLine(depth int) int {
	move := e.ProbePvTable()
	count := 0

	for move != c.NOMOVE && count < depth {
		if e.MoveExists(move) {
			e.MakeMove(move)
			e.Board.PvArray[count] = move
			count++
		} else {
			break
		}
		move = e.ProbePvTable()
	}

	for e.Board.Ply > 0 {
		e.TakeMove()
	}

	return count
}

func (e *Engine) ProbeHashEntry(move, score *int, alpha, beta, depth int) bool {
	idx := e.Board.PosKey % uint64(e.Board.HashTable.NumEntries)

	if e.Board.HashTable.HashEntries[idx].PosKey == e.Board.PosKey {
		*move = e.Board.HashTable.HashEntries[idx].Move
		if e.Board.HashTable.HashEntries[idx].Depth >= depth {
			e.Board.HashTable.Hit++
			*score = e.Board.HashTable.HashEntries[idx].Score
			if *score > c.MATE {
				*score -= e.Board.Ply
			} else if *score < -c.MATE {
				*score += e.Board.Ply
			}

			switch e.Board.HashTable.HashEntries[idx].Flags {
			case int(c.HFALPHA):
				if *score <= alpha {
					*score = alpha
					return true
				}
			case int(c.HFBETA):
				if *score >= beta {
					*score = beta
					return true
				}
			case int(c.HFEXACT):
				return true
			}
		}
	}

	return false
}

func (e *Engine) StoreHashEntry(move, score, flags, depth int) {
	idx := e.Board.PosKey % uint64(e.Board.HashTable.NumEntries)

	if e.Board.HashTable.HashEntries[idx].PosKey == 0 {
		e.Board.HashTable.NewWrite++
	} else {
		e.Board.HashTable.OverWrite++
	}

	if score > c.MATE {
		score += e.Board.Ply
	} else if score < -c.MATE {
		score -= e.Board.Ply
	}

	e.Board.HashTable.HashEntries[idx].Move = move
	e.Board.HashTable.HashEntries[idx].Score = score
	e.Board.HashTable.HashEntries[idx].Flags = flags
	e.Board.HashTable.HashEntries[idx].Depth = depth
	e.Board.HashTable.HashEntries[idx].PosKey = e.Board.PosKey
}
