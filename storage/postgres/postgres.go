package postgres

import (
	"context"
	"fmt"
	"github.com/Ramazon1227/go-rest-api-starter/config"
	"github.com/Ramazon1227/go-rest-api-starter/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db   *pgxpool.Pool
	user storage.UserRepoImpl
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:   pool,
		user: NewUserRepo(pool),
	}, nil
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) User() storage.UserRepoImpl {
	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}

	return s.user
}
