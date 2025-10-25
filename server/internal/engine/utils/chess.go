package utils

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
)

func Fr2Sq(f c.File, r c.Rank) c.Square {
	return (c.Square(f) + 21) + (c.Square(r) * 10)
}

func Sq64(sq120 c.Square) int {
	return c.Sq120ToSq64[sq120]
}

func Sq120(sq64 int) c.Square {
	return c.Sq64ToSq120[sq64]
}

func FromSq(m int) c.Square {
	return c.Square(m & 0x7F)
}

func ToSq(m int) c.Square {
	return c.Square((m >> 7) & 0x7F)
}

func Captured(m int) c.Piece {
	return c.Piece((m >> 14) & 0xF)
}

func Promoted(m int) c.Piece {
	return c.Piece((m >> 20) & 0xF)
}
