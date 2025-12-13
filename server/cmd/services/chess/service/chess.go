package service

import (
	t "github.com/jsbento/chess-server-v4/cmd/services/chess/types"
	e "github.com/jsbento/chess-server-v4/internal/engine"
	o "github.com/jsbento/chess-server-v4/internal/engine/opening-book"
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
	openingBook, err := o.NewOpeningBook(engine)
	if err != nil {
		return nil, err
	}
	engine.WithOpeningBook(openingBook)
	if err := engine.ParseFEN(req.Fen); err != nil {
		return nil, err
	}

	info := &eT.SearchInfo{}
	move, openingName := engine.ParseGo(req.ToGoCmd(), info)

	resp := &t.MoveResp{
		Move: move,
	}
	if openingName != "" {
		resp.OpeningName = &openingName
	}

	return resp, nil
}
