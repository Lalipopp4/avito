package redis

import (
	"context"
)

type Repository interface {
	CountUsers(ctx context.Context) (int, error)
	CheckSegments(ctx context.Context, segment string) (bool, error)
	GetElementById(ctx context.Context, id int) (string, error)
}
