package behavior

import (
	"github.com/jsbento/chess-server-v4/pkg/db"
)

type Store struct {
	pg *db.Postgres
}

func NewStore(dsn string) (*Store, error) {
	pg, err := db.NewPostgres(dsn)
	if err != nil {
		return nil, err
	}
	return &Store{pg: pg}, nil
}

func (s *Store) Close() error {
	return s.pg.Close()
}
