package engine

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	t "github.com/jsbento/chess-server-v4/internal/engine/types"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

type Engine struct {
	Board      *t.Board
	SetMask    [64]uint64
	ClearMask  [64]uint64
	PieceKeys  [13][120]uint64
	SideKey    uint64
	EnPasKey   uint64
	CastleKeys [16]uint64
}

func NewEngine() (e *Engine) {
	e = &Engine{
		Board:      t.NewBoard(),
		SetMask:    [64]uint64{},
		ClearMask:  [64]uint64{},
		PieceKeys:  [13][120]uint64{},
		SideKey:    0,
		EnPasKey:   0,
		CastleKeys: [16]uint64{},
	}
	e.InitBitmasks()
	e.InitHashKeys()
	return e
}

func (e *Engine) InitBitmasks() {
	for i := 0; i < 64; i++ {
		e.SetMask[i] = 0
		e.ClearMask[i] = 0
	}

	for i := 0; i < 64; i++ {
		e.SetMask[i] |= (uint64(1) << uint64(i))
		e.ClearMask[i] = ^e.SetMask[i]
	}
}

func (e *Engine) SetBit(bb *uint64, sq int) {
	*bb |= e.SetMask[sq]
}

func (e *Engine) ClearBit(bb *uint64, sq int) {
	*bb &= e.ClearMask[sq]
}

func (e *Engine) ParseMove(m string) int {
	if m[1] > '8' || m[1] < '1' {
		return c.NOMOVE
	}
	if m[3] > '8' || m[3] < '1' {
		return c.NOMOVE
	}
	if m[0] > 'h' || m[0] < 'a' {
		return c.NOMOVE
	}
	if m[2] > 'h' || m[2] < 'a' {
		return c.NOMOVE
	}

	from := utils.Fr2Sq(c.File((m[0] - 'a')), c.Rank(m[1]-'1'))
	to := utils.Fr2Sq(c.File(m[2]-'a'), c.Rank(m[3]-'1'))

	list := t.NewMoveList()
	e.GenerateAllMoves(list)

	for i := 0; i < list.Count; i++ {
		move := list.Moves[i].Move
		if utils.FromSq(move) == from && utils.ToSq(move) == to {
			promPce := int(utils.Promoted(move))
			if promPce != int(c.EMPTY) {
				if utils.IsRQ(promPce) && !utils.IsBQ(promPce) && m[4] == 'r' {
					return move
				} else if !utils.IsRQ(promPce) && utils.IsBQ(promPce) && m[4] == 'b' {
					return move
				} else if utils.IsRQ(promPce) && utils.IsBQ(promPce) && m[4] == 'q' {
					return move
				} else if utils.IsKn(promPce) && m[4] == 'n' {
					return move
				}
				continue
			}
			return move
		}
	}

	return c.NOMOVE
}
