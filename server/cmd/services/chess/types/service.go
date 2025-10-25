package types

import "fmt"

type EvalPosReq struct {
	Fen string `json:"fen"`
}

type PosScoreResp struct {
	Score int `json:"score"`
}

type SearchPosReq struct {
	Fen      string `json:"fen"`
	MoveTime int    `json:"moveTime"`
	Depth    int    `json:"depth"`
}

type MoveResp struct {
	Move string `json:"move"`
}

func (r *SearchPosReq) ToGoCmd() string {
	out := "go"

	if r.Depth > 0 {
		out += fmt.Sprintf(" depth %d", r.Depth)
	}
	if r.MoveTime > 0 {
		out += fmt.Sprintf(" movetime %d", r.MoveTime)
	}

	return out
}
