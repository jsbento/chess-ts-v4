package engine

import (
	"math"

	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	d "github.com/jsbento/chess-server-v4/internal/engine/data"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

var PawnTable [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 5, 5, 10, 10, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KnightTable [64]int = [64]int{
	0, -10, 0, 0, 0, 0, -10, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 5, 0,
	5, 10, 15, 20, 20, 15, 10, 5,
	5, 10, 10, 20, 20, 10, 10, 5,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}
var BishopTable [64]int = [64]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var RookTable [64]int = [64]int{
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	25, 25, 25, 25, 25, 25, 25, 25,
	0, 0, 5, 10, 10, 5, 0, 0,
}

var KingEndgame [64]int = [64]int{
	-50, -10, 0, 0, 0, 0, -10, -50,
	-10, 0, 10, 10, 10, 10, 0, -10,
	0, 10, 15, 15, 15, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 15, 15, 15, 10, 0,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-50, -10, 0, 0, 0, 0, -10, -50,
}

var KingTable [64]int = [64]int{
	0, 5, 5, -10, -10, 0, 10, 5,
	-30, -30, -30, -30, -30, -30, -30, -30,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
}

var _Mirror64 [64]int = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

const (
	PawnIsolated      int = -10
	RookOpenFile      int = 10
	RookSemiOpenFile  int = 5
	QueenOpenFile     int = 5
	QueenSemiOpenFile int = 3
	BishopPair        int = 30
)

var PawnPassed [8]int = [8]int{0, 5, 10, 20, 35, 60, 100, 200}
var EndgameMaterial int = d.PieceVal[c.WR] + 2*d.PieceVal[c.WN] + 2*d.PieceVal[c.WP]

func Mirror64(sq int) int {
	return _Mirror64[sq]
}

func (e *Engine) MaterialDraw() bool {
	if e.Board.PceNum[c.WR] == 0 && e.Board.PceNum[c.BR] == 0 && e.Board.PceNum[c.WQ] == 0 &&
		e.Board.PceNum[c.BQ] == 0 {
		if e.Board.PceNum[c.BB] == 0 && e.Board.PceNum[c.WB] == 0 {
			if e.Board.PceNum[c.WN] < 3 && e.Board.PceNum[c.BN] < 3 {
				return true
			}
		} else if e.Board.PceNum[c.WN] == 0 && e.Board.PceNum[c.BN] == 0 {
			if math.Abs(float64(e.Board.PceNum[c.WB]-e.Board.PceNum[c.BB])) < 2 {
				return true
			}
		} else if (e.Board.PceNum[c.WN] < 3 && e.Board.PceNum[c.WB] == 0) || (e.Board.PceNum[c.WB] == 1 && e.Board.PceNum[c.WN] == 0) {
			if (e.Board.PceNum[c.BN] < 3 && e.Board.PceNum[c.BB] == 0) || (e.Board.PceNum[c.BB] == 1 && e.Board.PceNum[c.BN] == 0) {
				return true
			}
		}
	} else if e.Board.PceNum[c.WQ] == 0 && e.Board.PceNum[c.BQ] == 0 {
		if e.Board.PceNum[c.WR] == 1 && e.Board.PceNum[c.BR] == 1 {
			if (e.Board.PceNum[c.WN]+e.Board.PceNum[c.WB]) < 2 && (e.Board.PceNum[c.BN]+e.Board.PceNum[c.BB]) < 2 {
				return true
			}
		} else if e.Board.PceNum[c.WR] == 1 && e.Board.PceNum[c.BR] == 0 {
			if (e.Board.PceNum[c.WN]+e.Board.PceNum[c.WB] == 0) && (((e.Board.PceNum[c.BN] + e.Board.PceNum[c.BB]) == 1) || ((e.Board.PceNum[c.BN] + e.Board.PceNum[c.BB]) == 2)) {
				return true
			}
		} else if e.Board.PceNum[c.WR] == 0 && e.Board.PceNum[c.BR] == 1 {
			if (e.Board.PceNum[c.BN]+e.Board.PceNum[c.BB] == 0) && (((e.Board.PceNum[c.WN] + e.Board.PceNum[c.WB]) == 1) || ((e.Board.PceNum[c.WN] + e.Board.PceNum[c.WB]) == 2)) {
				return true
			}
		}
	}

	return false
}

func (e *Engine) EvalPosition() int {
	if e.Board.PceNum[c.WP] == 0 && e.Board.PceNum[c.BP] == 0 && e.MaterialDraw() {
		return 0
	}

	score := e.Board.Material[c.WHITE] - e.Board.Material[c.BLACK]

	pce := c.WP
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score += PawnTable[utils.Sq64(sq)]

		if c.IsolatedMask[utils.Sq64(sq)]&e.Board.Pawns[c.WHITE] == 0 {
			score += PawnIsolated
		}

		if c.WhitePassedMask[utils.Sq64(sq)]&e.Board.Pawns[c.BLACK] == 0 {
			score += PawnPassed[c.RanksBrd[sq]]
		}
	}

	pce = c.BP
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score -= PawnTable[Mirror64(utils.Sq64(sq))]

		if c.IsolatedMask[utils.Sq64(sq)]&e.Board.Pawns[c.BLACK] == 0 {
			score -= PawnIsolated
		}

		if c.BlackPassedMask[utils.Sq64(sq)]&e.Board.Pawns[c.WHITE] == 0 {
			score -= PawnPassed[7-c.RanksBrd[sq]]
		}
	}

	pce = c.WN
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score += KnightTable[utils.Sq64(sq)]
	}

	pce = c.BN
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score -= KnightTable[Mirror64(utils.Sq64(sq))]
	}

	pce = c.WB
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score += BishopTable[utils.Sq64(sq)]
	}

	pce = c.BB
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score -= BishopTable[Mirror64(utils.Sq64(sq))]
	}

	pce = c.WR
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score += RookTable[utils.Sq64(sq)]

		if e.Board.Pawns[c.BOTH]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score += RookOpenFile
		} else if e.Board.Pawns[c.WHITE]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score += RookSemiOpenFile
		}
	}

	pce = c.BR
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]
		score -= RookTable[Mirror64(utils.Sq64(sq))]

		if e.Board.Pawns[c.BOTH]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score -= RookOpenFile
		} else if e.Board.Pawns[c.BLACK]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score -= RookSemiOpenFile
		}
	}

	pce = c.WQ
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]

		if e.Board.Pawns[c.BOTH]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score += QueenOpenFile
		} else if e.Board.Pawns[c.WHITE]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score += QueenSemiOpenFile
		}
	}

	pce = c.BQ
	for i := 0; i < e.Board.PceNum[pce]; i++ {
		sq := e.Board.Plist[pce][i]

		if e.Board.Pawns[c.BOTH]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score -= QueenOpenFile
		} else if e.Board.Pawns[c.BLACK]&c.FileBBMask[c.FilesBrd[sq]] == 0 {
			score -= QueenSemiOpenFile
		}
	}

	pce = c.WK
	sq := e.Board.Plist[pce][0]
	if e.Board.Material[c.BLACK] <= EndgameMaterial {
		score += KingEndgame[utils.Sq64(sq)]
	} else {
		score += KingTable[utils.Sq64(sq)]
	}

	pce = c.BK
	sq = e.Board.Plist[pce][0]
	if e.Board.Material[c.WHITE] <= EndgameMaterial {
		score -= KingEndgame[Mirror64(utils.Sq64(sq))]
	} else {
		score -= KingTable[Mirror64(utils.Sq64(sq))]
	}

	if e.Board.PceNum[c.WB] >= 2 {
		score += BishopPair
	}
	if e.Board.PceNum[c.BB] >= 2 {
		score -= BishopPair
	}

	if e.Board.Side == c.WHITE {
		return score
	} else {
		return -score
	}
}
