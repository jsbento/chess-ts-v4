package engine

import (
	"fmt"

	t "github.com/jsbento/chess-server-v4/internal/engine/types"
	eUtils "github.com/jsbento/chess-server-v4/internal/engine/utils"
	"github.com/jsbento/chess-server-v4/pkg/utils"
)

var leafNodes uint64

func (e *Engine) Perft(depth int) {
	if depth == 0 {
		leafNodes++
		return
	}

	moveList := t.NewMoveList()
	e.GenerateAllMoves(moveList)

	for i := 0; i < moveList.Count; i++ {
		if !e.MakeMove(moveList.Moves[i].Move) {
			continue
		}
		e.Perft(depth - 1)
		e.TakeMove()
	}
}

func (e *Engine) PerftTest(depth int) {
	fmt.Printf("\nStarting Perft Test To Depth:%d\n", depth)
	leafNodes = 0

	start := utils.GetTimeMs()
	moveList := t.NewMoveList()
	e.GenerateAllMoves(moveList)

	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i].Move
		if !e.MakeMove(move) {
			continue
		}
		cumNodes := leafNodes
		e.Perft(depth - 1)
		e.TakeMove()
		oldNodes := leafNodes - cumNodes
		fmt.Printf("move %d : %s : %d\n", i+1, eUtils.PrintMove(move), oldNodes)
	}

	fmt.Printf("\nTest Complete : %d nodes visited in %dms\n", leafNodes, utils.GetTimeMs()-start)
}
