package segment

import (
	"context"

	"github.com/Lalipopp4/test_api/internal/models"
)

type SegmentService interface {
	AddSegment(ctx context.Context, segment *models.Segment) error
	DeleteSegment(ctx context.Context, segment *models.Segment) error
	GetHistoryByDate(ctx context.Context, date *models.Date) (string, error)
	GetCSV(ctx context.Context, date string) ([]byte, error)
}
