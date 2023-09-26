package segment

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/Lalipopp4/test_api/internal/models"
)

const (
	HISTORYFILEPATH = "pkg/files/history/history"
)

func (s *segmentService) AddSegment(ctx context.Context, segment *models.Segment) error {

	// checking if segment doesn't exist
	exists, err := s.redisRepo.CheckSegments(ctx, segment.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("Segment " + segment.Name + " already exists.")
	}

	// chosing users for this segment
	n, err := s.redisRepo.CountUsers(ctx)
	if err != nil {
		return err
	}
	users := make([]int, n)
	usersMap := make(map[int]struct{})
	for i := 0; i < n/100*segment.Perc; {
		fmt.Println(i)
		v := rand.Intn(n + 1)
		if _, ok := usersMap[v]; !ok {
			users = append(users, v)
			usersMap[v] = struct{}{}
			i++
		}
	}

	// adding segment
	err = s.psqlRepo.AddSegment(ctx, segment, users)
	if err != nil {
		return err
	}
	return nil
}

func (s *segmentService) DeleteSegment(ctx context.Context, segment *models.Segment) error {

	//checking if segment exists
	exists, err := s.redisRepo.CheckSegments(ctx, segment.Name)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Segment " + segment.Name + " doesn't exist.")
	}

	// getting segment id
	segmentId, err := s.psqlRepo.GetSegmentIdByName(ctx, segment.Name)
	if err != nil {
		return err
	}

	// getting users in this segment
	users, err := s.psqlRepo.Get(ctx, segmentId, false, -1)
	if err != nil {
		return err
	}
	// deleting segment
	err = s.psqlRepo.DeleteSegment(ctx, segment, users)
	if err != nil {
		return err
	}

	return nil
}

func (s *segmentService) GetHistoryByDate(ctx context.Context, date *models.Date) (string, error) {

	// getting history by date
	history, err := s.psqlRepo.GetHistoryByDate(ctx, date.Date)
	if err != nil {
		return "", err
	}
	if len(history) == 0 {
		return "", errors.New("No records on " + date.Date + ".")
	}
	f, err := os.Create(HISTORYFILEPATH + date.Date + ".csv")
	defer f.Close()

	if err != nil {
		return "", err
	}
	for _, val := range history {
		if _, err := f.WriteString(val[0] + ";" + val[1] + ";" + val[2] + ";" + val[3] + "\n"); err != nil {
			return "", err
		}
	}
	return "/csvhistory?date=" + date.Date, nil
}

func (s *segmentService) GetCSV(ctx context.Context, date string) ([]byte, error) {
	f, err := os.Open(HISTORYFILEPATH + date + ".csv")
	if err != nil {
		return nil, err
	}
	temp := make([]byte, 100)
	data := []byte{}
	for {
		_, err := f.Read(temp)
		if err == io.EOF {
			break
		}
		data = append(data, temp...)
	}
	return data, nil
}

func (s *segmentService) CheckTTL() {
	for {
		s.psqlRepo.DeleteTTL(context.Background(), time.Now().Format(time.DateOnly))
		time.Sleep(time.Hour * 24)
	}
}
