package engine

import (
	"errors"
	"fmt"

	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	d "github.com/jsbento/chess-server-v4/internal/engine/data"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

func (e *Engine) ResetBoard() {
	for i := 0; i < c.BRD_SQ_NUM; i++ {
		e.Board.Pieces[i] = c.Piece(c.OFFBOARD)
	}

	for i := 0; i < 64; i++ {
		e.Board.Pieces[c.Sq64ToSq120[i]] = c.Piece(c.EMPTY)
	}

	for i := 0; i < 2; i++ {
		e.Board.BigPce[i] = 0
		e.Board.MajPce[i] = 0
		e.Board.MinPce[i] = 0
		e.Board.Material[i] = 0
	}

	for i := 0; i < 3; i++ {
		e.Board.Pawns[i] = uint64(0)
	}

	for i := 0; i < 13; i++ {
		e.Board.PceNum[i] = 0
	}

	e.Board.KingSq[c.WHITE] = int(c.NO_SQ)
	e.Board.KingSq[c.BLACK] = int(c.NO_SQ)

	e.Board.Side = c.BOTH
	e.Board.EnPas = c.NO_SQ
	e.Board.FiftyMove = 0

	e.Board.Ply = 0
	e.Board.HisPly = 0

	e.Board.CastlePerm = 0

	e.Board.PosKey = uint64(0)
}

func (e *Engine) PrintBoard() {
	board := e.Board

	fmt.Println("Game Board:")

	for rank := c.RANK_8; rank >= c.RANK_1; rank-- {
		fmt.Printf("%d  ", rank+1)
		for file := c.FILE_A; file <= c.FILE_H; file++ {
			sq := utils.Fr2Sq(file, rank)
			piece := board.Pieces[sq]
			fmt.Printf("%3c", d.PceChar[piece])
		}
		fmt.Println()
	}

	fmt.Println("\n       ")
	for file := c.FILE_A; file <= c.FILE_H; file++ {
		fmt.Printf("%3c", 'a'+int(file))
	}
	fmt.Println()
	fmt.Printf("side:%c\n", d.SideChar[board.Side])
	fmt.Printf("enPas:%d\n", board.EnPas)

	castleStr := ""
	if board.CastlePerm&c.WKCA != 0 {
		castleStr += "K"
	} else {
		castleStr += "-"
	}
	if board.CastlePerm&c.WQCA != 0 {
		castleStr += "Q"
	} else {
		castleStr += "-"
	}
	if board.CastlePerm&c.BKCA != 0 {
		castleStr += "k"
	} else {
		castleStr += "-"
	}
	if board.CastlePerm&c.BQCA != 0 {
		castleStr += "q"
	} else {
		castleStr += "-"
	}

	fmt.Printf("castle:%s\n", castleStr)
	fmt.Printf("PosKey:%x\n", board.PosKey)
}

func (e *Engine) ParseFEN(fen string) error {
	e.ResetBoard()

	rank := c.RANK_8
	file := c.FILE_A
	piece := c.Piece(0)
	count := 0
	i := 0
	sq64 := 0
	sq120 := c.Square(0)
	fenIdx := 0

	for rank >= c.RANK_1 && fenIdx < len(fen) {
		count = 1
		switch fen[fenIdx] {
		case 'p':
			piece = c.BP
		case 'r':
			piece = c.BR
		case 'n':
			piece = c.BN
		case 'b':
			piece = c.BB
		case 'k':
			piece = c.BK
		case 'q':
			piece = c.BQ
		case 'P':
			piece = c.WP
		case 'R':
			piece = c.WR
		case 'N':
			piece = c.WN
		case 'B':
			piece = c.WB
		case 'K':
			piece = c.WK
		case 'Q':
			piece = c.WQ
		case '1', '2', '3', '4', '5', '6', '7', '8':
			piece = c.Piece(c.EMPTY)
			count = int(fen[fenIdx]) - int('0')
		case '/', ' ':
			rank--
			file = c.FILE_A
			fenIdx++
			continue
		default:
			return errors.New("invalid FEN string")
		}

		for i = 0; i < count; i++ {
			sq64 = int(rank)*8 + int(file)
			sq120 = c.Sq64ToSq120[sq64]
			if piece != c.Piece(c.EMPTY) {
				e.Board.Pieces[sq120] = piece
			}
			file++
		}
		fenIdx++
	}

	if fen[fenIdx] == 'w' {
		e.Board.Side = c.WHITE
	} else {
		e.Board.Side = c.BLACK
	}
	fenIdx += 2

	for i = 0; i < 4; i++ {
		if fen[fenIdx] == ' ' {
			break
		}
		switch fen[fenIdx] {
		case 'K':
			e.Board.CastlePerm |= c.WKCA
		case 'Q':
			e.Board.CastlePerm |= c.WQCA
		case 'k':
			e.Board.CastlePerm |= c.BKCA
		case 'q':
			e.Board.CastlePerm |= c.BQCA
		default:
		}
		fenIdx++
	}
	fenIdx++

	if fen[fenIdx] != '-' {
		file = c.File(int(fen[fenIdx]) - int('a'))
		rank = c.Rank(int(fen[fenIdx+1]) - int('1'))

		e.Board.EnPas = utils.Fr2Sq(file, rank)
	}

	e.GeneratePosKey()
	e.UpdateListsMaterial()

	return nil
}

func (e *Engine) UpdateListsMaterial() {
	for i := 0; i < c.BRD_SQ_NUM; i++ {
		pce := e.Board.Pieces[i]
		if pce != c.Piece(c.OFFBOARD) && pce != c.Piece(c.EMPTY) {
			color := d.PieceCol[pce]

			if d.PieceBig[pce] {
				e.Board.BigPce[color]++
			}
			if d.PieceMin[pce] {
				e.Board.MinPce[color]++
			}
			if d.PieceMaj[pce] {
				e.Board.MajPce[color]++
			}

			e.Board.Material[color] += int(d.PieceVal[pce])
			e.Board.Plist[pce][e.Board.PceNum[pce]] = c.Square(i)
			e.Board.PceNum[pce]++

			if pce == c.WK {
				e.Board.KingSq[c.WHITE] = i
			}
			if pce == c.BK {
				e.Board.KingSq[c.BLACK] = i
			}

			if pce == c.WP {
				e.SetBit(&e.Board.Pawns[c.WHITE], utils.Sq64(c.Square(i)))
				e.SetBit(&e.Board.Pawns[c.BOTH], utils.Sq64(c.Square(i)))
			} else if pce == c.BP {
				e.SetBit(&e.Board.Pawns[c.BLACK], utils.Sq64(c.Square(i)))
				e.SetBit(&e.Board.Pawns[c.BOTH], utils.Sq64(c.Square(i)))
			}
		}
	}
}
