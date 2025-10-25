package utils

import (
	"fmt"

	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	d "github.com/jsbento/chess-server-v4/internal/engine/data"
)

func IsKn(pce int) bool {
	return d.PieceKnight[pce]
}

func IsKi(pce int) bool {
	return d.PieceKing[pce]
}

func IsRQ(pce int) bool {
	return d.PieceRookQueen[pce]
}

func IsBQ(pce int) bool {
	return d.PieceBishopQueen[pce]
}

func PrintMove(m int) string {
	ff := c.FilesBrd[FromSq(m)]
	rf := c.RanksBrd[FromSq(m)]
	ft := c.FilesBrd[ToSq(m)]
	rt := c.RanksBrd[ToSq(m)]

	promPce := Promoted(m)
	if promPce != c.Piece(c.EMPTY) {
		promChar := 'q'
		if IsKn(int(promPce)) {
			promChar = 'n'
		} else if IsRQ(int(promPce)) && !IsBQ(int(promPce)) {
			promChar = 'r'
		} else if !IsRQ(int(promPce)) && IsBQ(int(promPce)) {
			promChar = 'b'
		}
		return fmt.Sprintf("%c%c%c%c%c", 'a'+ff, '1'+rf, 'a'+ft, '1'+rt, promChar)
	} else {
		return fmt.Sprintf("%c%c%c%c", 'a'+ff, '1'+rf, 'a'+ft, '1'+rt)
	}
}

func PrintSquare(sq c.Square) string {
	file := c.FilesBrd[sq]
	rank := c.RanksBrd[sq]
	return fmt.Sprintf("%c%c", 'a'+file, '1'+rank)
}
