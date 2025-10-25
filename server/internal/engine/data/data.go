package data

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
)

const (
	PceChar  string = ".PNBRQKpnbrqk"
	SideChar string = "wb-"
	RankChar string = "12345678"
	FileChar string = "abcdefgh"
)

var PieceBig [13]bool = [13]bool{
	false, false, true, true, true, true, true,
	false, true, true, true, true, true,
}

var PieceMaj [13]bool = [13]bool{
	false, false, false, false, true, true, true,
	false, false, false, true, true, true,
}

var PieceMin [13]bool = [13]bool{
	false, false, true, true, false, false, false,
	false, true, true, false, false, false,
}

var PieceVal [13]int = [13]int{
	0, 100, 325, 325, 550, 1000, 50000,
	100, 325, 325, 550, 1000, 50000,
}

var PieceCol [13]c.Side = [13]c.Side{
	c.BOTH, c.WHITE, c.WHITE, c.WHITE, c.WHITE, c.WHITE, c.WHITE,
	c.BLACK, c.BLACK, c.BLACK, c.BLACK, c.BLACK, c.BLACK,
}

var PieceKnight [13]bool = [13]bool{
	false, false, true, false, false, false, false,
	false, true, false, false, false, false,
}

var PieceKing [13]bool = [13]bool{
	false, false, false, false, false, false, true,
	false, false, false, false, false, true,
}

var PieceRookQueen [13]bool = [13]bool{
	false, false, false, false, true, true, false,
	false, false, false, true, true, false,
}

var PieceBishopQueen [13]bool = [13]bool{
	false, false, false, true, false, true, false,
	false, false, true, false, true, false,
}

var PiecePawn [13]bool = [13]bool{
	false, true, false, false, false, false, false,
	true, false, false, false, false, false,
}

var PieceSlides [13]bool = [13]bool{
	false, false, false, true, true, true, false,
	false, false, true, true, true, false,
}

var VictimScore [13]int = [13]int{
	0, 100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600,
}

var MvvLvaScores [13][13]int
