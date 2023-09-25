package redis

import (
	"context"

	"github.com/Lalipopp4/test_api/internal/config"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	cur *redis.Client
}

func New(cfg config.Redis) (Repository, error) {
	db := &repository{redis.NewClient(&redis.Options{
		Addr:     cfg.Host + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})}
	if _, err := db.cur.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return db, nil
}
