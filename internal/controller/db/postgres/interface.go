package postgres

import (
	"context"

	"github.com/Lalipopp4/test_api/internal/models"
)

type Repository interface {
	AddUser(ctx context.Context, user *models.UserRequest, segmentsToAdd []int, segmentsToDelete []int) error
	AddSegment(ctx context.Context, segment *models.Segment, users []int) error
	DeleteSegment(ctx context.Context, segment *models.Segment, users []int) error
	GetHistoryByDate(ctx context.Context, date string) ([][4]string, error)
	Get(ctx context.Context, id int, filter bool, extra int) ([]int, error)
	GetSegmentIdByName(ctx context.Context, name string) (int, error)
	DeleteTTL(ctx context.Context, date string) error
}
