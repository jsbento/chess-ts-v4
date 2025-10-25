package bitboards

import (
	c "github.com/jsbento/chess-server/pkg/constants"
	"github.com/jsbento/chess-server/pkg/utils"
)

func PopBit(bb *uint64) int {
	b := uint64(*bb ^ (*bb - 1))
	fold := uint32((b & 0xffffffff) ^ (b >> 32))
	*bb &= (*bb - 1)
	return c.BitTable[(fold*uint32(0x783a9b23))>>26]
}

func CountBits(b uint64) int {
	var count int
	for b != 0 {
		b &= b - 1
		count++
	}
	return count
}

func PrintBitBoard(bb uint64) {
	shift := uint64(1)

	for rank := c.RANK_8; rank >= c.RANK_1; rank-- {
		for file := c.FILE_A; file <= c.FILE_H; file++ {
			sq := utils.Fr2Sq(file, rank)
			sq64 := c.Sq120ToSq64[sq]
			if (shift<<sq64)&bb != 0 {
				print("X")
			} else {
				print("-")
			}
		}
		println()
	}
}
