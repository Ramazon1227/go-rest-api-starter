package mongo

import (
	"context"

	"github.com/Ramazon1227/go-rest-api-starter/config"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store struct {
	client *mongo.Client
	db     *mongo.Database
	user   storage.UserRepoImpl
}

func NewMongo(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.MongoDatabase)

	return &Store{
		client: client,
		db:     db,
		user:   NewUserRepo(db),
	}, nil
}

func (s *Store) CloseDB() {
	s.client.Disconnect(context.Background()) //nolint:errcheck
}

func (s *Store) User() storage.UserRepoImpl {
	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}
	return s.user
}
