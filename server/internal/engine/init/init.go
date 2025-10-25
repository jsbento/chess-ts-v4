package init

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	d "github.com/jsbento/chess-server-v4/internal/engine/data"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

func AllInit() {
	InitSq120To64()
	InitFilesRanksBrd()
	InitMvvLva()
	InitEvalMasks()
}

func InitSq120To64() {
	for i := 0; i < c.BRD_SQ_NUM; i++ {
		c.Sq120ToSq64[i] = 65
	}

	for i := 0; i < 64; i++ {
		c.Sq64ToSq120[i] = 120
	}

	sq64 := 0
	for rank := c.RANK_1; rank <= c.RANK_8; rank++ {
		for file := c.FILE_A; file <= c.FILE_H; file++ {
			sq := utils.Fr2Sq(file, rank)
			c.Sq64ToSq120[sq64] = sq
			c.Sq120ToSq64[int(sq)] = sq64
			sq64++
		}
	}
}

func InitFilesRanksBrd() {
	for i := 0; i < c.BRD_SQ_NUM; i++ {
		c.FilesBrd[i] = int(c.OFFBOARD)
		c.RanksBrd[i] = int(c.OFFBOARD)
	}

	for rank := c.RANK_1; rank <= c.RANK_8; rank++ {
		for file := c.FILE_A; file <= c.FILE_H; file++ {
			sq := utils.Fr2Sq(file, rank)
			c.FilesBrd[sq] = int(file)
			c.RanksBrd[sq] = int(rank)
		}
	}
}

func InitMvvLva() {
	for attacker := c.WP; attacker < c.BK; attacker++ {
		for victim := c.WP; victim < c.BK; victim++ {
			d.MvvLvaScores[victim][attacker] = d.VictimScore[victim] + 6 - (d.VictimScore[attacker] / 100)
		}
	}
}

func InitEvalMasks() {
	for i := 0; i < 8; i++ {
		c.FileBBMask[i] = uint64(0)
		c.RankBBMask[i] = uint64(0)
	}

	for rank := c.RANK_8; rank >= c.RANK_1; rank-- {
		for file := c.FILE_A; file <= c.FILE_H; file++ {
			sq := uint64(rank)*8 + uint64(file)
			c.FileBBMask[file] |= uint64(1) << sq
			c.RankBBMask[rank] |= uint64(1) << sq
		}
	}

	for sq := 0; sq < 64; sq++ {
		c.IsolatedMask[sq] = uint64(0)
		c.WhitePassedMask[sq] = uint64(0)
		c.BlackPassedMask[sq] = uint64(0)
	}

	for sq := 0; sq < 64; sq++ {
		tSq := sq + 8

		for tSq < 64 {
			c.WhitePassedMask[sq] |= uint64(1) << uint64(tSq)
			tSq += 8
		}

		tSq = sq - 8
		for tSq >= 0 {
			c.BlackPassedMask[sq] |= uint64(1) << uint64(tSq)
			tSq -= 8
		}

		if c.FilesBrd[utils.Sq120(sq)] > int(c.FILE_A) {
			c.IsolatedMask[sq] |= c.FileBBMask[c.FilesBrd[utils.Sq120(sq)]-1]

			tSq = sq + 7
			for tSq < 64 {
				c.WhitePassedMask[sq] |= uint64(1) << uint64(tSq)
				tSq += 8
			}

			tSq = sq - 9
			for tSq >= 0 {
				c.BlackPassedMask[sq] |= uint64(1) << uint64(tSq)
				tSq -= 8
			}
		}

		if c.FilesBrd[utils.Sq120(sq)] < int(c.FILE_H) {
			c.IsolatedMask[sq] |= c.FileBBMask[c.FilesBrd[utils.Sq120(sq)]+1]

			tSq = sq + 9
			for tSq < 64 {
				c.WhitePassedMask[sq] |= uint64(1) << uint64(tSq)
				tSq += 8
			}

			tSq = sq - 7
			for tSq >= 0 {
				c.BlackPassedMask[sq] |= uint64(1) << uint64(tSq)
				tSq -= 8
			}
		}
	}
}
