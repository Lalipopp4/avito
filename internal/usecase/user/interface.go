package user

import (
	"context"

	"github.com/Lalipopp4/test_api/internal/models"
)

type UserService interface {
	AddUser(ctx context.Context, user *models.UserRequest) error
	GetSegmentsByUser(ctx context.Context, user *models.User) ([]string, error)
}
