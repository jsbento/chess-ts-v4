package engine

import (
	"math"

	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	t "github.com/jsbento/chess-server-v4/internal/engine/types"
	eUtils "github.com/jsbento/chess-server-v4/internal/engine/utils"
	"github.com/jsbento/chess-server-v4/pkg/utils"
)

func (e *Engine) PickNextMove(moveNum int, list *t.MoveList) {
	bestScore := 0
	bestNum := moveNum
	for i := moveNum; i < list.Count; i++ {
		if list.Moves[i].Score > bestScore {
			bestScore = list.Moves[i].Score
			bestNum = i
		}
	}
	temp := list.Moves[moveNum]
	list.Moves[moveNum] = list.Moves[bestNum]
	list.Moves[bestNum] = temp
}

func (e *Engine) IsRepetition() bool {
	for i := e.Board.HisPly - e.Board.FiftyMove; i < e.Board.HisPly-1; i++ {
		if e.Board.PosKey == e.Board.History[i].PosKey {
			return true
		}
	}
	return false
}

func (e *Engine) CheckUp(info *t.SearchInfo) {
	if info.TimeSet && utils.GetTimeMs() > info.StopTime {
		info.Stopped = true
	}
}

func (e *Engine) ClearForSearch(info *t.SearchInfo) {
	for i := 0; i < 13; i++ {
		for j := 0; j < c.BRD_SQ_NUM; j++ {
			e.Board.SearchHistory[i][j] = 0
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < c.MAX_DEPTH; j++ {
			e.Board.SearchKillers[i][j] = 0
		}
	}

	e.Board.HashTable.OverWrite = 0
	e.Board.HashTable.Hit = 0
	e.Board.HashTable.Cut = 0
	e.Board.Ply = 0

	info.Stopped = false
	info.Nodes = 0
	info.Fh = 0.0
	info.Fhf = 0.0
}

func (e *Engine) Quiescence(alpha, beta int, info *t.SearchInfo) int {
	if info.Nodes&2047 == 0 {
		e.CheckUp(info)
	}

	info.Nodes++

	if e.IsRepetition() || e.Board.FiftyMove >= 100 {
		return 0
	}

	if e.Board.Ply > c.MAX_DEPTH-1 {
		return e.EvalPosition()
	}

	score := e.EvalPosition()

	if score >= beta {
		return beta
	}

	if score > alpha {
		alpha = score
	}

	list := t.NewMoveList()
	e.GenerateAllCaptures(list)

	legal := 0
	score = -c.INFINITE

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		e.PickNextMove(moveNum, list)
		if !e.MakeMove(list.Moves[moveNum].Move) {
			continue
		}

		legal++
		score = -e.Quiescence(-beta, -alpha, info)
		e.TakeMove()

		if score > alpha {
			if score >= beta {
				if legal == 1 {
					info.Fhf++
				}
				info.Fh++

				return beta
			}
			alpha = score
		}
	}

	return alpha
}

func (e *Engine) AlphaBeta(alpha, beta, depth int, info *t.SearchInfo, doNull bool) int {
	if depth == 0 {
		return e.Quiescence(alpha, beta, info)
	}

	if info.Nodes&2047 == 0 {
		e.CheckUp(info)
	}

	info.Nodes++

	if e.IsRepetition() || e.Board.FiftyMove >= 100 && e.Board.Ply != 0 {
		return 0
	}

	if e.Board.Ply > c.MAX_DEPTH-1 {
		return e.EvalPosition()
	}

	inCheck := e.IsSqAttacked(c.Square(e.Board.KingSq[e.Board.Side]), e.Board.Side^1)
	if inCheck {
		depth++
	}

	score := -c.INFINITE
	pvMove := c.NOMOVE

	if e.ProbeHashEntry(&pvMove, &score, alpha, beta, depth) {
		e.Board.HashTable.Cut++
		return score
	}

	if doNull && !inCheck && e.Board.Ply != 0 && e.Board.BigPce[e.Board.Side] > 1 && depth >= 4 {
		e.MakeNullMove()
		score = -e.AlphaBeta(-beta, -beta+1, depth-4, info, false)
		e.TakeNullMove()

		if info.Stopped {
			return 0
		}

		if score >= beta && int(math.Abs(float64(score))) < c.MATE {
			info.NullCut++
			return beta
		}
	}

	list := t.NewMoveList()
	e.GenerateAllMoves(list)

	legal := 0
	oldAlpha := alpha
	bestMove := c.NOMOVE
	bestScore := -c.INFINITE
	score = -c.INFINITE

	if pvMove != c.NOMOVE {
		for moveNum := 0; moveNum < list.Count; moveNum++ {
			if list.Moves[moveNum].Move == pvMove {
				list.Moves[moveNum].Score = 2000000
				break
			}
		}
	}

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		e.PickNextMove(moveNum, list)
		if !e.MakeMove(list.Moves[moveNum].Move) {
			continue
		}

		legal++
		score = -e.AlphaBeta(-beta, -alpha, depth-1, info, true)
		e.TakeMove()

		if info.Stopped {
			return 0
		}
		if score > bestScore {
			bestScore = score
			bestMove = list.Moves[moveNum].Move
			if score > alpha {
				if score >= beta {
					if legal == 1 {
						info.Fhf++
					}
					info.Fh++

					if list.Moves[moveNum].Move&int(c.MFLAGCAP) == 0 {
						e.Board.SearchKillers[1][e.Board.Ply] = e.Board.SearchKillers[0][e.Board.Ply]
						e.Board.SearchKillers[0][e.Board.Ply] = list.Moves[moveNum].Move
					}

					e.StoreHashEntry(bestMove, beta, int(c.HFBETA), depth)

					return beta
				}
				alpha = score
				bestMove = list.Moves[moveNum].Move

				if list.Moves[moveNum].Move&int(c.MFLAGCAP) == 0 {
					e.Board.SearchHistory[e.Board.Pieces[eUtils.FromSq(bestMove)]][eUtils.ToSq(bestMove)] += depth
				}
			}
		}
	}

	if legal == 0 {
		if inCheck {
			return -c.MATE + e.Board.Ply
		} else {
			return 0
		}
	}

	if alpha != oldAlpha {
		e.StoreHashEntry(bestMove, bestScore, int(c.HFEXACT), depth)
	} else {
		e.StoreHashEntry(bestMove, alpha, int(c.HFALPHA), depth)
	}

	return alpha
}

func (e *Engine) SearchPosition(info *t.SearchInfo) (string, string) {
	bestMove := c.NOMOVE

	if e.OpeningBook != nil {
		if e.OpeningBook.InBook(e.Board.PosKey) {
			moves := e.OpeningBook.GetMoves(e.Board.PosKey)
			if len(moves) > 0 {
				mv := moves[utils.RandRange(0, len(moves)-1)]
				if e.ParseMove(mv) != c.NOMOVE {
					return mv, e.OpeningBook.GetOpeningName(e.Board.PosKey)[0]
				}
			}
		}
	}

	e.ClearForSearch(info)
	for currDepth := 1; currDepth <= info.Depth; currDepth++ {
		e.AlphaBeta(-c.INFINITE, c.INFINITE, currDepth, info, true)

		if info.Stopped {
			break
		}

		e.GetPvLine(currDepth)
		bestMove = e.Board.PvArray[0]
	}

	return eUtils.PrintMove(bestMove), ""
}
