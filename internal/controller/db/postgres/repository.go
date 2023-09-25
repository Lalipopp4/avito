package postgres

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/Lalipopp4/test_api/internal/models"
)

func (r repository) AddSegment(ctx context.Context, segment *models.Segment, users []int) error {
	tx, err := r.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var (
		txErr error
		wg    sync.WaitGroup
		mu    sync.Mutex
	)

	wg.Add(1 + 2*len(users))
	go func() {
		mu.Lock()
		_, err = r.cur.ExecContext(ctx, "INSERT INTO segments (name) VALUES ($1)", segment.Name)
		mu.Unlock()
		if err != nil {
			txErr = err
		}
		wg.Done()
	}()
	for _, val := range users {
		go func() {
			mu.Lock()
			_, err := tx.ExecContext(ctx, "INSERT INTO user_segments (user_id, segment_id) VALUES ($1, $2)", val, segment.Id)
			mu.Unlock()
			if err != nil {
				txErr = err
			}
		}()
		go func() {
			mu.Lock()
			_, err = tx.ExecContext(ctx, `INSERT INTO history (user_id, segment_id, operation, "time")
			 VALUES ($1, $2, $3, $4)`, val, segment.Id, true, time.Now())
			mu.Unlock()
			if err != nil {
				txErr = err
			}

		}()

	}

	wg.Wait()
	if txErr != nil {
		tx.Rollback()
		return txErr
	}

	tx.Commit()
	return nil
}

func (r *repository) DeleteSegment(ctx context.Context, segment *models.Segment, users []int) error {
	tx, err := r.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var (
		txErr error
		wg    sync.WaitGroup
		mu    sync.Mutex
	)

	wg.Add(len(users) + 2)
	go func() {
		defer wg.Done()
		mu.Lock()
		_, err = tx.ExecContext(ctx, "DELETE FROM segments WHERE id = $1", segment.Id)
		mu.Unlock()
		if err != nil {
			txErr = err
		}
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		_, err = tx.ExecContext(ctx, "DELETE FROM user_segments WHERE segment_id = $1", segment.Id)
		mu.Unlock()
		if err != nil {
			txErr = err
		}
	}()

	for _, val := range users {
		go func() {
			defer wg.Done()
			mu.Lock()
			_, err = tx.ExecContext(ctx, `INSERT INTO history (user_id, segment_id, operation, "time")
        	 VALUES ($1, $2, $3, $4)`, val, segment.Id, false, time.Now())
			mu.Unlock()
			if err != nil {
				txErr = err
			}
		}()
	}

	wg.Wait()
	if txErr != nil {
		tx.Rollback()
		return txErr
	}
	tx.Commit()
	return nil
}

func (r *repository) GetHistoryByDate(ctx context.Context, date string) ([][4]string, error) {
	rows, err := r.cur.QueryContext(ctx, `SELECT * FROM history WHERE 
	EXTRACT(YEAR FROM "time") = $1 AND EXTRACT(MONTH FROM "time") = $2`, date[:4], date[5:])

	if err != nil {
		return nil, err
	}
	var (
		res               [][4]string
		userId, segmentId int
		operation         bool
		timestamp, opS    string
		wg                sync.WaitGroup
	)
	for rows.Next() {
		wg.Add(1)
		go func() {
			rows.Scan(&userId, &segmentId, &operation, &date)
			if operation {
				opS = "Adding"
			} else {
				opS = "Deleting"
			}
			res = append(res, [4]string{strconv.Itoa(userId), strconv.Itoa(segmentId), opS, timestamp})
			wg.Done()
		}()
	}
	wg.Wait()
	return res, nil
}

func (r *repository) AddUser(ctx context.Context, user *models.UserRequest, segmentsToAdd []int, segmentsToDelete []int) error {
	tx, err := r.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var (
		wg sync.WaitGroup
	)
	wg.Add(len(segmentsToAdd) + len(segmentsToDelete))
	for i, segment := range user.SegmentsToAdd {
		go func() {
			defer wg.Done()
			r.AddSegment(ctx, &models.Segment{Name: segment, Id: segmentsToAdd[i]}, []int{user.Id})

		}()
	}

	for i, segment := range user.SegmentsToDelete {
		go func() {
			defer wg.Done()
			r.DeleteSegment(ctx, &models.Segment{Name: segment, Id: segmentsToDelete[i]}, []int{user.Id})
		}()
	}

	tx.Commit()
	return nil
}

func (r *repository) GetSegmentIdByName(ctx context.Context, name string) (int, error) {
	rows := r.cur.QueryRowContext(ctx, `SELECT id FROM segments WHERE name=$1`, name)
	var (
		segmentId int
	)
	rows.Scan(&segmentId)
	return segmentId, nil
}

func (r *repository) Get(ctx context.Context, id int, filter bool, extra int) ([]int, error) {
	var (
		cond, goal string
		extraCond  string
	)
	if extra != -1 {
		extraCond = " AND id_user = " + strconv.Itoa(extra)
	}
	switch filter {
	case false:
		cond = "id_segment"
		goal = "id_user"
	case true:
		cond = "id_user"
		goal = "id_segment"
	}
	rows, err := r.cur.QueryContext(ctx, "SELECT "+goal+" FROM user_segments WHERE "+cond+" =$1"+extraCond, id)
	if err != nil {
		return nil, err
	}
	var (
		res []int
		wg  sync.WaitGroup
	)
	for rows.Next() {
		go func() {
			var temp int
			rows.Scan(&temp)
			res = append(res, temp)
			wg.Done()
		}()
	}

	wg.Wait()
	return res, nil
}
