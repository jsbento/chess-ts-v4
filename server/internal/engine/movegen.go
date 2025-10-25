package engine

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	d "github.com/jsbento/chess-server-v4/internal/engine/data"
	t "github.com/jsbento/chess-server-v4/internal/engine/types"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

var LoopSlidePiece [8]c.Piece = [8]c.Piece{
	c.WB, c.WR, c.WQ, 0, c.BQ, c.BR, c.BB, 0,
}
var LoopSlideIndex [2]int = [2]int{0, 4}

var LoopNonSlidePiece [6]c.Piece = [6]c.Piece{
	c.WN, c.WK, 0, c.BK, c.BN, 0,
}
var LoopNonSlideIndex [2]int = [2]int{0, 3}

var PieceDir [13][8]int = [13][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
}

var NumDir [13]int = [13]int{
	0, 0, 8, 4, 4, 8, 8, 0, 8, 4, 4, 8, 8,
}

func Move(from c.Square, to c.Square, captured c.Piece, promoted c.Piece, flag c.MoveFlag) int {
	return int(from) | (int(to) << 7) | (int(captured) << 14) | (int(promoted) << 20) | int(flag)
}

func SqOffboard(sq c.Square) bool {
	return c.FilesBrd[sq] == int(c.OFFBOARD)
}

func (e *Engine) MoveExists(move int) bool {
	moveList := t.NewMoveList()
	e.GenerateAllMoves(moveList)

	for i := 0; i < moveList.Count; i++ {
		if !e.MakeMove(moveList.Moves[i].Move) {
			continue
		}
		e.TakeMove()
		if move == moveList.Moves[i].Move {
			return true
		}
	}

	return false
}

func (e *Engine) AddQuietMove(move int, list *t.MoveList) {
	list.Moves[list.Count].Move = move

	if e.Board.SearchKillers[0][e.Board.Ply] == move {
		list.Moves[list.Count].Score = 900000
	} else if e.Board.SearchKillers[1][e.Board.Ply] == move {
		list.Moves[list.Count].Score = 800000
	} else {
		list.Moves[list.Count].Score = e.Board.SearchHistory[e.Board.Pieces[utils.FromSq(move)]][utils.ToSq(move)]
	}

	list.Count++
}

func (e *Engine) AddCaptureMove(move int, list *t.MoveList) {
	list.Moves[list.Count].Move = move
	list.Moves[list.Count].Score = d.MvvLvaScores[utils.Captured(move)][e.Board.Pieces[utils.FromSq(move)]] + 1000000
	list.Count++
}

func (e *Engine) AddEnPassantMove(move int, list *t.MoveList) {
	list.Moves[list.Count].Move = move
	list.Moves[list.Count].Score = 1000105
	list.Count++
}

func (e *Engine) AddWhitePawnCapMove(from, to c.Square, cap c.Piece, list *t.MoveList) {
	if c.RanksBrd[from] == int(c.RANK_7) {
		e.AddCaptureMove(Move(from, to, cap, c.WQ, 0), list)
		e.AddCaptureMove(Move(from, to, cap, c.WR, 0), list)
		e.AddCaptureMove(Move(from, to, cap, c.WB, 0), list)
		e.AddCaptureMove(Move(from, to, cap, c.WN, 0), list)
	} else {
		e.AddCaptureMove(Move(from, to, cap, c.EMPTY, 0), list)
	}
}

func (e *Engine) AddWhitePawnMove(from, to c.Square, list *t.MoveList) {
	if c.RanksBrd[from] == int(c.RANK_7) {
		e.AddQuietMove(Move(from, to, c.EMPTY, c.WQ, 0), list)
		e.AddQuietMove(Move(from, to, c.EMPTY, c.WR, 0), list)
		e.AddQuietMove(Move(from, to, c.EMPTY, c.WB, 0), list)
		e.AddQuietMove(Move(from, to, c.EMPTY, c.WN, 0), list)
	} else {
		e.AddQuietMove(Move(from, to, c.EMPTY, c.EMPTY, 0), list)
	}
}

func (e *Engine) AddBlackPawnCapMove(from, to c.Square, cap c.Piece, list *t.MoveList) {
	if c.RanksBrd[from] == int(c.RANK_2) {
		e.AddCaptureMove(Move(from, to, cap, c.BQ, 0), list)
		e.AddCaptureMove(Move(from, to, cap, c.BR, 0), list)
		e.AddCaptureMove(Move(from, to, cap, c.BB, 0), list)
		e.AddCaptureMove(Move(from, to, cap, c.BN, 0), list)
	} else {
		e.AddCaptureMove(Move(from, to, cap, c.EMPTY, 0), list)
	}
}

func (e *Engine) AddBlackPawnMove(from, to c.Square, list *t.MoveList) {
	if c.RanksBrd[from] == int(c.RANK_2) {
		e.AddQuietMove(Move(from, to, c.EMPTY, c.BQ, 0), list)
		e.AddQuietMove(Move(from, to, c.EMPTY, c.BR, 0), list)
		e.AddQuietMove(Move(from, to, c.EMPTY, c.BB, 0), list)
		e.AddQuietMove(Move(from, to, c.EMPTY, c.BN, 0), list)
	} else {
		e.AddQuietMove(Move(from, to, c.EMPTY, c.EMPTY, 0), list)
	}
}

func (e *Engine) GenerateAllCaptures(list *t.MoveList) {
	list.Count = 0

	if e.Board.Side == c.WHITE {
		for pceNum := 0; pceNum < e.Board.PceNum[c.WP]; pceNum++ {
			sq := e.Board.Plist[c.WP][pceNum]

			if !SqOffboard(sq+9) && d.PieceCol[int(e.Board.Pieces[sq+9])] == c.BLACK {
				e.AddWhitePawnCapMove(sq, sq+9, e.Board.Pieces[sq+9], list)
			}
			if !SqOffboard(sq+11) && d.PieceCol[int(e.Board.Pieces[sq+11])] == c.BLACK {
				e.AddWhitePawnCapMove(sq, sq+11, e.Board.Pieces[sq+11], list)
			}

			if e.Board.EnPas != c.NO_SQ {
				if sq+9 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq+9, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
				if sq+11 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq+11, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
			}
		}
	} else {
		for pceNum := 0; pceNum < e.Board.PceNum[c.BP]; pceNum++ {
			sq := e.Board.Plist[c.BP][pceNum]

			if !SqOffboard(sq-9) && d.PieceCol[int(e.Board.Pieces[sq-9])] == c.WHITE {
				e.AddBlackPawnCapMove(sq, sq-9, e.Board.Pieces[sq-9], list)
			}
			if !SqOffboard(sq-11) && d.PieceCol[int(e.Board.Pieces[sq-11])] == c.WHITE {
				e.AddBlackPawnCapMove(sq, sq-11, e.Board.Pieces[sq-11], list)
			}

			if e.Board.EnPas != c.NO_SQ {
				if sq-9 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq-9, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
				if sq-11 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq-11, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
			}
		}
	}

	pceIdx := LoopSlideIndex[e.Board.Side]
	pce := LoopSlidePiece[pceIdx]
	pceIdx++

	for pce != 0 {
		for pceNum := 0; pceNum < e.Board.PceNum[pce]; pceNum++ {
			sq := e.Board.Plist[pce][pceNum]
			for i := 0; i < NumDir[pce]; i++ {
				dir := PieceDir[pce][i]
				tSq := sq + c.Square(dir)

				for !SqOffboard(tSq) {
					if e.Board.Pieces[tSq] != c.Piece(c.EMPTY) {
						if d.PieceCol[int(e.Board.Pieces[tSq])] == e.Board.Side^1 {
							e.AddCaptureMove(Move(sq, tSq, e.Board.Pieces[tSq], c.EMPTY, 0), list)
						}
						break
					}
					tSq += c.Square(dir)
				}
			}
		}

		pce = LoopSlidePiece[pceIdx]
		pceIdx++
	}

	pceIdx = LoopNonSlideIndex[e.Board.Side]
	pce = LoopNonSlidePiece[pceIdx]
	pceIdx++

	for pce != 0 {
		for pceNum := 0; pceNum < e.Board.PceNum[pce]; pceNum++ {
			sq := e.Board.Plist[pce][pceNum]
			for i := 0; i < NumDir[pce]; i++ {
				dir := PieceDir[pce][i]
				tSq := sq + c.Square(dir)

				if SqOffboard(tSq) {
					continue
				}

				if e.Board.Pieces[tSq] != c.Piece(c.EMPTY) {
					if d.PieceCol[int(e.Board.Pieces[tSq])] == e.Board.Side^1 {
						e.AddCaptureMove(Move(sq, tSq, e.Board.Pieces[tSq], c.EMPTY, 0), list)
					}
					continue
				}
			}
		}

		pce = LoopNonSlidePiece[pceIdx]
		pceIdx++
	}
}

func (e *Engine) GenerateAllMoves(list *t.MoveList) {
	list.Count = 0

	if e.Board.Side == c.WHITE {
		for pceNum := 0; pceNum < e.Board.PceNum[c.WP]; pceNum++ {
			sq := e.Board.Plist[c.WP][pceNum]

			if e.Board.Pieces[sq+10] == c.Piece(c.EMPTY) {
				e.AddWhitePawnMove(sq, sq+10, list)
				if c.RanksBrd[sq] == int(c.RANK_2) && e.Board.Pieces[sq+20] == c.Piece(c.EMPTY) {
					e.AddQuietMove(Move(sq, sq+20, c.EMPTY, c.EMPTY, c.MFLAGPS), list)
				}
			}

			if !SqOffboard(sq+9) && d.PieceCol[int(e.Board.Pieces[sq+9])] == c.BLACK {
				e.AddWhitePawnCapMove(sq, sq+9, e.Board.Pieces[sq+9], list)
			}
			if !SqOffboard(sq+11) && d.PieceCol[int(e.Board.Pieces[sq+11])] == c.BLACK {
				e.AddWhitePawnCapMove(sq, sq+11, e.Board.Pieces[sq+11], list)
			}

			if e.Board.EnPas != c.NO_SQ {
				if sq+9 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq+9, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
				if sq+11 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq+11, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
			}
		}

		if e.Board.CastlePerm&c.WKCA != 0 {
			if e.Board.Pieces[c.Square(c.F1)] == c.Piece(c.EMPTY) &&
				e.Board.Pieces[c.Square(c.G1)] == c.Piece(c.EMPTY) {
				if !e.IsSqAttacked(c.F1, c.BLACK) && !e.IsSqAttacked(c.E1, c.BLACK) {
					e.AddQuietMove(Move(c.E1, c.G1, c.EMPTY, c.EMPTY, c.MFLAGCA), list)
				}
			}
		}

		if e.Board.CastlePerm&c.WQCA != 0 {
			if e.Board.Pieces[c.Square(c.D1)] == c.Piece(c.EMPTY) &&
				e.Board.Pieces[c.Square(c.C1)] == c.Piece(c.EMPTY) &&
				e.Board.Pieces[c.Square(c.B1)] == c.Piece(c.EMPTY) {
				if !e.IsSqAttacked(c.D1, c.BLACK) && !e.IsSqAttacked(c.E1, c.BLACK) {
					e.AddQuietMove(Move(c.E1, c.C1, c.EMPTY, c.EMPTY, c.MFLAGCA), list)
				}
			}
		}
	} else {
		for pceNum := 0; pceNum < e.Board.PceNum[c.BP]; pceNum++ {
			sq := e.Board.Plist[c.BP][pceNum]

			if e.Board.Pieces[sq-10] == c.Piece(c.EMPTY) {
				e.AddBlackPawnMove(sq, sq-10, list)
				if c.RanksBrd[sq] == int(c.RANK_7) && e.Board.Pieces[sq-20] == c.Piece(c.EMPTY) {
					e.AddQuietMove(Move(sq, sq-20, c.EMPTY, c.EMPTY, c.MFLAGPS), list)
				}
			}

			if !SqOffboard(sq-9) && d.PieceCol[int(e.Board.Pieces[sq-9])] == c.WHITE {
				e.AddBlackPawnCapMove(sq, sq-9, e.Board.Pieces[sq-9], list)
			}
			if !SqOffboard(sq-11) && d.PieceCol[int(e.Board.Pieces[sq-11])] == c.WHITE {
				e.AddBlackPawnCapMove(sq, sq-11, e.Board.Pieces[sq-11], list)
			}

			if e.Board.EnPas != c.NO_SQ {
				if sq-9 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq-9, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
				if sq-11 == e.Board.EnPas {
					e.AddEnPassantMove(Move(sq, sq-11, c.EMPTY, c.EMPTY, c.MFLAGEP), list)
				}
			}
		}

		if e.Board.CastlePerm&c.BKCA != 0 {
			if e.Board.Pieces[c.Square(c.F8)] == c.Piece(c.EMPTY) && e.Board.Pieces[c.Square(c.G8)] == c.Piece(c.EMPTY) {
				if !e.IsSqAttacked(c.F8, c.WHITE) && !e.IsSqAttacked(c.E8, c.WHITE) {
					e.AddQuietMove(Move(c.E8, c.G8, c.EMPTY, c.EMPTY, c.MFLAGCA), list)
				}
			}
		}

		if e.Board.CastlePerm&c.BQCA != 0 {
			if e.Board.Pieces[c.Square(c.D8)] == c.Piece(c.EMPTY) && e.Board.Pieces[c.Square(c.C8)] == c.Piece(c.EMPTY) && e.Board.Pieces[c.Square(c.B8)] == c.Piece(c.EMPTY) {
				if !e.IsSqAttacked(c.Square(c.D8), c.WHITE) && !e.IsSqAttacked(c.Square(c.E8), c.WHITE) {
					e.AddQuietMove(Move(c.E8, c.C8, c.EMPTY, c.EMPTY, c.MFLAGCA), list)
				}
			}
		}
	}

	pceIdx := LoopSlideIndex[e.Board.Side]
	pce := LoopSlidePiece[pceIdx]
	pceIdx++

	for pce != 0 {
		for pceNum := 0; pceNum < e.Board.PceNum[pce]; pceNum++ {
			sq := e.Board.Plist[pce][pceNum]
			for i := 0; i < NumDir[pce]; i++ {
				dir := PieceDir[pce][i]
				tSq := sq + c.Square(dir)

				for !SqOffboard(tSq) {
					if e.Board.Pieces[tSq] != c.Piece(c.EMPTY) {
						if d.PieceCol[int(e.Board.Pieces[tSq])] == e.Board.Side^1 {
							e.AddCaptureMove(Move(sq, tSq, e.Board.Pieces[tSq], c.EMPTY, 0), list)
						}
						break
					}
					e.AddQuietMove(Move(sq, tSq, c.EMPTY, c.EMPTY, 0), list)
					tSq += c.Square(dir)
				}
			}
		}

		pce = LoopSlidePiece[pceIdx]
		pceIdx++
	}

	pceIdx = LoopNonSlideIndex[e.Board.Side]
	pce = LoopNonSlidePiece[pceIdx]
	pceIdx++

	for pce != 0 {
		for pceNum := 0; pceNum < e.Board.PceNum[pce]; pceNum++ {
			sq := e.Board.Plist[pce][pceNum]
			for i := 0; i < NumDir[pce]; i++ {
				dir := PieceDir[pce][i]
				tSq := sq + c.Square(dir)

				if SqOffboard(tSq) {
					continue
				}

				if e.Board.Pieces[tSq] != c.Piece(c.EMPTY) {
					if d.PieceCol[int(e.Board.Pieces[tSq])] == e.Board.Side^1 {
						e.AddCaptureMove(Move(sq, tSq, e.Board.Pieces[tSq], c.EMPTY, 0), list)
					}
					continue
				}
				e.AddQuietMove(Move(sq, tSq, c.EMPTY, c.EMPTY, 0), list)
			}
		}

		pce = LoopNonSlidePiece[pceIdx]
		pceIdx++
	}
}
