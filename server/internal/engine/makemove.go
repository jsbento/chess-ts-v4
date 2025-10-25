package engine

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	d "github.com/jsbento/chess-server-v4/internal/engine/data"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

var CastlePerm [120]int = [120]int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

func (e *Engine) HashPce(pce c.Piece, sq c.Square) {
	e.Board.PosKey ^= e.PieceKeys[pce][sq]
}

func (e *Engine) HashCastle() {
	e.Board.PosKey ^= e.CastleKeys[e.Board.CastlePerm]
}

func (e *Engine) HashSide() {
	e.Board.PosKey ^= e.SideKey
}

func (e *Engine) HashEnPas() {
	e.Board.PosKey ^= e.PieceKeys[c.EMPTY][e.Board.EnPas]
}

func (e *Engine) ClearPiece(sq c.Square) {
	pce := e.Board.Pieces[sq]
	col := d.PieceCol[pce]
	tPceNum := -1

	e.HashPce(pce, sq)

	e.Board.Pieces[sq] = c.EMPTY
	e.Board.Material[col] -= d.PieceVal[pce]

	if d.PieceBig[pce] {
		e.Board.BigPce[col]--
		if d.PieceMaj[pce] {
			e.Board.MajPce[col]--
		} else {
			e.Board.MinPce[col]--
		}
	} else {
		e.ClearBit(&e.Board.Pawns[col], utils.Sq64(sq))
		e.ClearBit(&e.Board.Pawns[c.BOTH], utils.Sq64(sq))
	}

	for i := 0; i < e.Board.PceNum[pce]; i++ {
		if e.Board.Plist[pce][i] == sq {
			tPceNum = i
			break
		}
	}

	e.Board.PceNum[pce]--
	e.Board.Plist[pce][tPceNum] = e.Board.Plist[pce][e.Board.PceNum[pce]]
}

func (e *Engine) AddPiece(sq c.Square, pce c.Piece) {
	col := d.PieceCol[pce]
	e.HashPce(pce, sq)

	e.Board.Pieces[sq] = pce

	if d.PieceBig[pce] {
		e.Board.BigPce[col]++
		if d.PieceMaj[pce] {
			e.Board.MajPce[col]++
		} else {
			e.Board.MinPce[col]++
		}
	} else {
		e.SetBit(&e.Board.Pawns[col], utils.Sq64(sq))
		e.SetBit(&e.Board.Pawns[c.BOTH], utils.Sq64(sq))
	}

	e.Board.Material[col] += d.PieceVal[pce]
	e.Board.Plist[pce][e.Board.PceNum[pce]] = sq
	e.Board.PceNum[pce]++
}

func (e *Engine) MovePiece(from, to c.Square) {
	pce := e.Board.Pieces[from]
	col := d.PieceCol[pce]

	e.HashPce(pce, from)
	e.Board.Pieces[from] = c.EMPTY

	e.HashPce(pce, to)
	e.Board.Pieces[to] = pce

	if !d.PieceBig[pce] {
		e.ClearBit(&e.Board.Pawns[col], utils.Sq64(from))
		e.ClearBit(&e.Board.Pawns[c.BOTH], utils.Sq64(from))
		e.SetBit(&e.Board.Pawns[col], utils.Sq64(to))
		e.SetBit(&e.Board.Pawns[c.BOTH], utils.Sq64(to))
	}

	for i := 0; i < e.Board.PceNum[pce]; i++ {
		if e.Board.Plist[pce][i] == from {
			e.Board.Plist[pce][i] = to
			break
		}
	}
}

func (e *Engine) MakeMove(move int) bool {
	from := utils.FromSq(move)
	to := utils.ToSq(move)
	side := e.Board.Side

	e.Board.History[e.Board.HisPly].PosKey = e.Board.PosKey

	if move&int(c.MFLAGEP) != 0 {
		if side == c.WHITE {
			e.ClearPiece(to - 10)
		} else {
			e.ClearPiece(to + 10)
		}
	} else if move&int(c.MFLAGCA) != 0 {
		switch to {
		case c.C1:
			e.MovePiece(c.A1, c.D1)
		case c.C8:
			e.MovePiece(c.A8, c.D8)
		case c.G1:
			e.MovePiece(c.H1, c.F1)
		case c.G8:
			e.MovePiece(c.H8, c.F8)
		}
	}

	if e.Board.EnPas != c.NO_SQ {
		e.HashEnPas()
	}
	e.HashCastle()

	e.Board.History[e.Board.HisPly].Move = move
	e.Board.History[e.Board.HisPly].FiftyMove = e.Board.FiftyMove
	e.Board.History[e.Board.HisPly].EnPas = e.Board.EnPas
	e.Board.History[e.Board.HisPly].CastlePerm = e.Board.CastlePerm

	e.Board.CastlePerm &= c.CastlePerm(CastlePerm[from])
	e.Board.CastlePerm &= c.CastlePerm(CastlePerm[to])
	e.Board.EnPas = c.NO_SQ

	e.HashCastle()

	cap := utils.Captured(move)
	e.Board.FiftyMove++

	if cap != c.EMPTY {
		e.ClearPiece(to)
		e.Board.FiftyMove = 0
	}

	e.Board.HisPly++
	e.Board.Ply++

	if d.PiecePawn[e.Board.Pieces[from]] {
		e.Board.FiftyMove = 0
		if move&int(c.MFLAGPS) != 0 {
			if side == c.WHITE {
				e.Board.EnPas = from + 10
			} else {
				e.Board.EnPas = from - 10
			}
			e.HashEnPas()
		}
	}

	e.MovePiece(from, to)

	promPce := utils.Promoted(move)
	if promPce != c.EMPTY {
		e.ClearPiece(to)
		e.AddPiece(to, promPce)
	}

	if d.PieceKing[e.Board.Pieces[to]] {
		e.Board.KingSq[side] = int(to)
	}

	e.Board.Side ^= 1
	e.HashSide()

	if e.IsSqAttacked(c.Square(e.Board.KingSq[side]), e.Board.Side) {
		e.TakeMove()
		return false
	}

	return true
}

func (e *Engine) TakeMove() {
	e.Board.HisPly--
	e.Board.Ply--

	move := e.Board.History[e.Board.HisPly].Move
	from := utils.FromSq(move)
	to := utils.ToSq(move)

	if e.Board.EnPas != c.NO_SQ {
		e.HashEnPas()
	}
	e.HashCastle()

	e.Board.CastlePerm = e.Board.History[e.Board.HisPly].CastlePerm
	e.Board.FiftyMove = e.Board.History[e.Board.HisPly].FiftyMove
	e.Board.EnPas = e.Board.History[e.Board.HisPly].EnPas

	if e.Board.EnPas != c.NO_SQ {
		e.HashEnPas()
	}
	e.HashCastle()

	e.Board.Side ^= 1
	e.HashSide()

	if move&int(c.MFLAGEP) != 0 {
		if e.Board.Side == c.WHITE {
			e.AddPiece(to-10, c.BP)
		} else {
			e.AddPiece(to+10, c.WP)
		}
	} else if move&int(c.MFLAGCA) != 0 {
		switch to {
		case c.C1:
			e.MovePiece(c.D1, c.A1)
		case c.C8:
			e.MovePiece(c.D8, c.A8)
		case c.G1:
			e.MovePiece(c.F1, c.H1)
		case c.G8:
			e.MovePiece(c.F8, c.H8)
		}
	}

	e.MovePiece(to, from)

	if d.PieceKing[e.Board.Pieces[from]] {
		e.Board.KingSq[e.Board.Side] = int(from)
	}

	cap := utils.Captured(move)
	if cap != c.EMPTY {
		e.AddPiece(to, cap)
	}

	if utils.Promoted(move) != c.EMPTY {
		e.ClearPiece(from)
		if d.PieceCol[utils.Promoted(move)] == c.WHITE {
			e.AddPiece(from, c.WP)
		} else {
			e.AddPiece(from, c.BP)
		}
	}
}

func (e *Engine) MakeNullMove() {
	e.Board.Ply++
	e.Board.History[e.Board.HisPly].PosKey = e.Board.PosKey

	if e.Board.EnPas != c.NO_SQ {
		e.HashEnPas()
	}

	e.Board.History[e.Board.HisPly].Move = c.NOMOVE
	e.Board.History[e.Board.HisPly].FiftyMove = e.Board.FiftyMove
	e.Board.History[e.Board.HisPly].EnPas = e.Board.EnPas
	e.Board.History[e.Board.HisPly].CastlePerm = e.Board.CastlePerm
	e.Board.EnPas = c.NO_SQ

	e.Board.Side ^= 1
	e.Board.HisPly++
	e.HashSide()
}

func (e *Engine) TakeNullMove() {
	e.Board.HisPly--
	e.Board.Ply--

	if e.Board.EnPas != c.NO_SQ {
		e.HashEnPas()
	}

	e.Board.CastlePerm = e.Board.History[e.Board.HisPly].CastlePerm
	e.Board.FiftyMove = e.Board.History[e.Board.HisPly].FiftyMove
	e.Board.EnPas = e.Board.History[e.Board.HisPly].EnPas

	if e.Board.EnPas != c.NO_SQ {
		e.HashEnPas()
	}

	e.Board.Side ^= 1
	e.Board.PosKey ^= e.SideKey
	e.HashSide()
}
