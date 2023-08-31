package db

import (
	"strconv"
	"time"

	"github.com/Lalipopp4/avito/internal/logger"
	"gorm.io/gorm"
)

// counts records
func (db dbHandler) Count(table, column, goal string) (int, error) {
	var (
		cond string
		c    int64
	)
	switch column {
	case "":
		cond = "1 = 1"
	default:
		cond = `"` + column + `" = '` + goal + `'`
	}
	db.cur.Table(table).Select("count(*)").Where(cond).Count(&c)
	return int(c), nil
}

// selects records
func (db dbHandler) Select(table, cond string, goal ...string) ([]string, error) {
	var res []string
	db.cur.Table(table).Select(goal).Where(cond).Find(&res)
	return res, nil
}

// inserts records
func (db dbHandler) Insert(table, columns string, data string) error {
	type req struct {
		Name string
	}
	d := req{data}
	err := db.cur.Table(table).Select(columns).Create(&d)
	if err.Error != nil {
		logger.Logger.Log(err)
		return err.Error
	}
	return nil
}

// provides transaction with adding or deleting segments
func (db dbHandler) ExecSegments(users []int, segments []string, operation bool, ttl string) error {
	return db.cur.Transaction(func(tx *gorm.DB) error {
		var search []string
		switch operation {
		case true:
			type req struct {
				UserID      int
				SegmentName string
				TTL         string
			}
			reqs := make([]req, len(users)+len(segments)-1)
			if len(users) > 1 {
				for i, val := range users {
					reqs[i] = req{val, segments[0], ttl}
				}
			} else {
				for i, val := range segments {
					reqs[i] = req{users[0], val, ttl}
				}
			}
			if ttl == "" {
				search = []string{"user_id", "segment_name"}
			} else {
				search = []string{"user_id", "segment_name", "ttl"}
			}
			if err := tx.Table("user_segments").Select(search).Create(reqs).Error; err != nil {
				tx.Rollback()
				logger.Logger.Log(err)
				return err
			}
		default:
			reqs := make([]segReq, len(users)+len(segments)-1)
			if len(users) > 1 {
				for i, val := range users {
					reqs[i] = segReq{UserID: val, SegmentName: segments[0]}
				}

				// deleting segment from list of segments
				if err := tx.Exec("DELETE FROM segments WHERE name = $1", segments[0]).Error; err != nil {
					logger.Logger.Log(err)
					tx.Rollback()

					return err
				}

			} else {
				for i, val := range segments {
					reqs[i] = segReq{UserID: users[0], SegmentName: val}
				}
			}
			// deleting records with exact segemnt or user
			for _, val := range reqs {
				if err := tx.Exec("DELETE FROM user_segments WHERE segment_name = $1 AND user_id = $2", val.SegmentName, val.UserID).Error; err != nil {
					tx.Rollback()
					logger.Logger.Log(err)
					return err
				}
			}

		}
		reqs := make([]historyReq, len(users)+len(segments)-1)
		if len(users) > 1 {
			for i, val := range users {
				reqs[i] = historyReq{val, segments[0], operation, time.Now()}
			}
		} else {
			for i, val := range segments {
				reqs[i] = historyReq{users[0], val, operation, time.Now()}
			}
		}

		// inserting history in segment history
		if err := tx.Table("segment_history").Select("user_id", "segment_name", "operation", "timestamp_req").Create(reqs).Error; err != nil {
			tx.Rollback()
			logger.Logger.Log(err)
			return err
		}
		return nil
	})
}

// select records with advanced condition
func (db dbHandler) SelectAdvanced(table, date string) ([][]string, error) {
	var res [][]string
	switch table {
	case "segment_history":
		var temp []historyReq

		// selecting history
		err := db.cur.Table("segment_history").Select("user_id", "segment_name", "operation", "timestamp_req").Where(`EXTRACT(YEAR FROM timestamp_req) = $1 AND EXTRACT(MONTH FROM timestamp_req) = $2`, date[:4], date[5:]).Find(&temp)
		if err.Error != nil {
			return nil, err.Error
		}
		var opS string
		res = make([][]string, len(temp))
		for i, val := range temp {
			if val.Operation {
				opS = "Adding"
			} else {
				opS = "Deleting"
			}
			res[i] = []string{strconv.Itoa(val.UserID), val.SegmentName, opS, val.Timestamp_req.String()}
		}
	default:
		var temp []segReq

		// selecting records with ttl
		err := db.cur.Table("user_segments").Select("user_id", "segment_name").Where(`EXTRACT(YEAR FROM ttl) = $1 AND EXTRACT(MONTH FROM ttl) = $2`, date[:4], date[5:7]).Find(&temp)
		if err.Error != nil {
			return nil, err.Error
		}
		res = make([][]string, len(temp))
		for i, val := range temp {
			res[i] = []string{strconv.Itoa(val.UserID), val.SegmentName}
		}
	}
	return res, nil
}
