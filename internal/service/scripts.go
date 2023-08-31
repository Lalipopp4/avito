package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Lalipopp4/avito/internal/db"
	"github.com/Lalipopp4/avito/internal/logger"
)

var (
	// map of temporary segments for user
	ttlUSers = make(map[string][]point)
)

// function to decode data in request
func decode(d interface{}, w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte(error.Error(err)))
		return err
	}
	return nil
}

func checkTTL(c chan int) {
	var wg sync.WaitGroup
	for {
		select {
		case <-c:
			return
		default:
		}
		res, err := db.DB.SelectAdvanced("user_segments", time.Now().Format(time.DateOnly))
		if err != nil {
			continue
		}
		if len(res) == 0 {
			time.Sleep(time.Second * 86400)
			continue
		}
		wg.Add(len(res))
		for _, val := range res {
			id, err := strconv.Atoi(val[0])
			if err != nil {
				continue
			}
			defer wg.Done()
			go db.DB.ExecSegments([]int{id}, []string{val[1]}, false, "")
		}
		wg.Wait()
		logger.Logger.Log("Deleted ttl users from their segments.")
	}
}
