package engine

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	"github.com/jsbento/chess-server-v4/pkg/utils"
)

func (e *Engine) InitHashKeys() {
	for i := range 13 {
		for j := 0; j < 120; j++ {
			e.PieceKeys[i][j] = utils.Rand64()
		}
	}
	e.SideKey = utils.Rand64()
	for i := range 16 {
		e.CastleKeys[i] = utils.Rand64()
	}
}

func (e *Engine) GeneratePosKey() {
	key := uint64(0)

	for sq := range c.BRD_SQ_NUM {
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
