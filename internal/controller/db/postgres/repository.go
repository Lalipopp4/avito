package postgres

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/Lalipopp4/test_api/internal/models"
	"golang.org/x/sync/errgroup"
)

func (r repository) AddSegment(ctx context.Context, segment *models.Segment, users []int) error {
	tx, err := r.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var (
		wg, errCtx = errgroup.WithContext(ctx)
		mu         sync.Mutex
	)

	wg.Go(func() error {
		if errCtx != nil {
			return errCtx.Err()
		}
		mu.Lock()
		_, err = r.cur.ExecContext(ctx, "INSERT INTO segments (name) VALUES ($1)", segment.Name)
		mu.Unlock()
		if err != nil {
			return err
		}
		return nil
	})

	for _, val := range users {
		val := val
		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}
			mu.Lock()
			_, err := tx.ExecContext(ctx, "INSERT INTO user_segments (user_id, segment_id) VALUES ($1, $2)", val, segment.Id)
			mu.Unlock()
			if err != nil {
				return err
			}
			return nil
		})

		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}
			mu.Lock()
			_, err = tx.ExecContext(ctx, `INSERT INTO history (user_id, segment_id, operation, "time")
			 VALUES ($1, $2, $3, $4)`, val, segment.Id, true, time.Now())
			mu.Unlock()
			if err != nil {
				return err
			}
			return nil
		})

	}

	if err := wg.Wait(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) DeleteSegment(ctx context.Context, segment *models.Segment, users []int) error {
	tx, err := r.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var (
		wg, errCtx = errgroup.WithContext(ctx)
		mu         sync.Mutex
	)

	wg.Go(func() error {
		if errCtx != nil {
			return errCtx.Err()
		}
		mu.Lock()
		_, err = tx.ExecContext(ctx, "DELETE FROM segments WHERE id = $1", segment.Id)
		mu.Unlock()
		if err != nil {
			return err
		}
		return nil
	})

	wg.Go(func() error {
		if errCtx != nil {
			return errCtx.Err()
		}
		mu.Lock()
		_, err = tx.ExecContext(ctx, "DELETE FROM user_segments WHERE segment_id = $1", segment.Id)
		mu.Unlock()
		if err != nil {
			return err
		}
		return nil
	})

	for _, val := range users {
		val := val
		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}
			mu.Lock()
			_, err = tx.ExecContext(ctx, `INSERT INTO history (user_id, segment_id, operation, "time")
        	 VALUES ($1, $2, $3, $4)`, val, segment.Id, false, time.Now())
			mu.Unlock()
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err = wg.Wait(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
		wg, errCtx = errgroup.WithContext(ctx)
	)

	for i, segment := range user.SegmentsToAdd {
		i, segment := i, segment
		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}
			return r.AddSegment(ctx, &models.Segment{Name: segment, Id: segmentsToAdd[i]}, []int{user.Id})
		})
	}

	for i, segment := range user.SegmentsToDelete {
		i, segment := i, segment
		wg.Go(func() error {
			if errCtx != nil {
				return errCtx.Err()
			}
			return r.DeleteSegment(ctx, &models.Segment{Name: segment, Id: segmentsToDelete[i]}, []int{user.Id})
		})
	}
	if err = wg.Wait(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
		wg.Add(1)
		go func() {
			defer wg.Done()
			var temp int
			rows.Scan(&temp)
			res = append(res, temp)
		}()
	}

	wg.Wait()
	return res, nil
}

func (r *repository) DeleteTTL(ctx context.Context, date string) error {

	rows, err := r.cur.QueryContext(ctx, "SELECT user_id, segment_id FROM user_segments WHERE ttl = $1", date)
	if err != nil {
		return err
	}
	var (
		wg, errCtx          = errgroup.WithContext(ctx)
		mu                  sync.Mutex
		delete              = [][2]int{}
		segment_id, user_id int
	)
	for rows.Next() {
		rows.Scan(&segment_id, &user_id)
		delete = append(delete, [2]int{segment_id, user_id})
	}

	if len(delete) == 0 {
		return nil
	}

	tx, err := r.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	wg.Go(func() error {
		if errCtx != nil {
			return errCtx.Err()
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM user_segments WHERE ttl = $1", date)
		if err != nil {
			return err
		}
		return nil
	})

	for _, val := range delete {
		segment_id, user_id := val[0], val[1]
		wg.Go(func() error {

			if errCtx != nil {
				return errCtx.Err()
			}
			mu.Lock()
			_, err = tx.ExecContext(ctx, `INSERT INTO history (user_id, segment_id, operation, "time")
             VALUES ($1, $2, $3, $4)`, user_id, segment_id, false, time.Now())
			mu.Unlock()
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err = wg.Wait(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
