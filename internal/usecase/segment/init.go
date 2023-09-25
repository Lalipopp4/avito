package segment

import (
	"github.com/Lalipopp4/test_api/internal/config"
	"github.com/Lalipopp4/test_api/internal/controller/db/postgres"
	"github.com/Lalipopp4/test_api/internal/controller/db/redis"
)

type segmentService struct {
	psqlRepo  postgres.Repository
	redisRepo redis.Repository
}

func New(cfg *config.Config) (SegmentService, error) {
	psqlRepo, err := postgres.New(cfg.Postgres)
	if err != nil {
		return nil, err
	}
	redisRepo, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, err
	}
	return &segmentService{
		psqlRepo,
		redisRepo,
	}, nil
}
