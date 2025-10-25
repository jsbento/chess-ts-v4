package engine

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	"github.com/jsbento/chess-server-v4/pkg/utils"
)

func (e *Engine) InitHashKeys() {
	for i := 0; i < 13; i++ {
		for j := 0; j < 120; j++ {
			e.PieceKeys[i][j] = utils.Rand64()
		}
	}
	e.SideKey = utils.Rand64()
	for i := 0; i < 16; i++ {
		e.CastleKeys[i] = utils.Rand64()
	}
}

func (e *Engine) GeneratePosKey() {
	key := uint64(0)

	for sq := 0; sq < c.BRD_SQ_NUM; sq++ {
		piece := e.Board.Pieces[sq]
		if int(piece) != int(c.NO_SQ) && int(piece) != int(c.EMPTY) &&
			int(piece) != int(c.OFFBOARD) {
			key ^= e.PieceKeys[piece][sq]
		}
	}

	if e.Board.Side == c.WHITE {
		key ^= e.SideKey
	}

	if e.Board.EnPas != c.NO_SQ {
		key ^= e.PieceKeys[c.EMPTY][e.Board.EnPas]
	}

	key ^= e.CastleKeys[e.Board.CastlePerm]

	e.Board.PosKey = key
}
