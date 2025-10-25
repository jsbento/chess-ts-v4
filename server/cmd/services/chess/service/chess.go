package service

import (
	t "github.com/jsbento/chess-server-v4/cmd/services/chess/types"
	e "github.com/jsbento/chess-server-v4/internal/engine"
	eT "github.com/jsbento/chess-server-v4/internal/engine/types"
)

func (s *ChessService) evaluatePosition(req *t.EvalPosReq) (*t.PosScoreResp, error) {
	engine := e.NewEngine()
	if err := engine.ParseFEN(req.Fen); err != nil {
		return nil, err
	}

	return &t.PosScoreResp{
		Score: engine.EvalPosition(),
	}, nil
}

func (s *ChessService) searchPosition(req *t.SearchPosReq) (*t.MoveResp, error) {
	engine := e.NewEngine()
	if err := engine.ParseFEN(req.Fen); err != nil {
		return nil, err
	}

	info := &eT.SearchInfo{}
	return &t.MoveResp{
		Move: engine.ParseGo(req.ToGoCmd(), info),
	}, nil
}
