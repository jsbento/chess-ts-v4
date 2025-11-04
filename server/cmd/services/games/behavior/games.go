package behavior

import (
	"github.com/google/uuid"
	t "github.com/jsbento/chess-server-v4/cmd/services/games/types"
)

func (s *Store) CreateGame(req *t.CreateGameReq) (*t.Game, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	game := &t.Game{
		Id:       uuid.New().String(),
		PlayerId: req.PlayerId,
		Moves:    req.Moves,
		Result:   req.Result,
	}
	if err := s.pg.GetDB().Create(game).Error; err != nil {
		return nil, err
	}
	return game, nil
}

func (s *Store) GetGame(id string) (*t.Game, error) {
	game := &t.Game{}
	if err := s.pg.GetDB().Model(&t.Game{}).Where("id = ?", id).First(game).Error; err != nil {
		return nil, err
	}
	return game, nil
}

func (s *Store) GetUserGames(playerId string) ([]*t.Game, error) {
	games := []*t.Game{}
	if err := s.pg.GetDB().Model(&t.Game{}).Where("player_id = ?", playerId).Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}
