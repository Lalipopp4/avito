package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/Lalipopp4/avito/internal/db"
	"github.com/Lalipopp4/avito/internal/logger"
)

// adding segments
func addSegment(w http.ResponseWriter, r *http.Request) {
	var segment segment
	err := decode(&segment, w, r)
	if err != nil {
		return
	}

	//checking if segment already exists
	n, err := db.DB.Count("segments", "name", segment.Name)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}

	if n > 0 || segment.Name == "" {
		logger.Logger.Log("Error: this segment is already exists or empty request.")
		w.Write([]byte("Error: this segment is already exists or empty request.\n"))
		return
	}

	// adding segment name into list of segments
	err = db.DB.Insert("segments", "name", segment.Name)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}

	// chosing users for this segment
	n, err = db.DB.Count("users", "", "")
	fmt.Println(segment.Perc)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	users := []int{}
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
	fmt.Println(len(users))

	// SQL transaction to insert segment and users
	err = db.DB.ExecSegments(users, []string{segment.Name}, true, "")
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}

	logger.Logger.Log(segment.Name + " added.")
	w.Write([]byte(segment.Name + " added.\n"))
}

// deleting segment
func deleteSegment(w http.ResponseWriter, r *http.Request) {
	var segment segment
	err := decode(&segment, w, r)
	if err != nil {
		return
	}
	//checking if segment doesn't exist
	n, err := db.DB.Count("segments", "name", segment.Name)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	if n == 0 {
		logger.Logger.Log("Error: this segment doesn't exist or empty request.\n")
		w.Write([]byte("Error: this segment doesn't exist or empty request.\n"))
		return
	}

	// SQL transaction to find users in this segment,
	// delete segment and insert info about deleting users from segment
	// ------
	usersS, err := db.DB.Select("user_segments", `segment_name = '`+segment.Name+`'`, "user_id")
	users := make([]int, len(usersS))
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}

	for i := range users {
		users[i], err = strconv.Atoi(usersS[i])
	}
	err = db.DB.ExecSegments(users, []string{segment.Name}, false, "")
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	// ------

	logger.Logger.Log(segment.Name + " deleted.")
	w.Write([]byte(segment.Name + " deleted."))
}

// adding user, adding and deleting segments for him
func addUser(w http.ResponseWriter, r *http.Request) {
	var request addUserRequest
	err := decode(&request, w, r)
	if err != nil {
		return
	}

	// checking if there is intersection of adding and deleting segments
	intersection := make(map[string]struct{})
	for _, segment := range request.SegmentsToAdd {
		intersection[segment] = struct{}{}
	}
	for _, val := range request.SegmentsToDelete {
		if _, ok := intersection[val]; ok {
			logger.Logger.Log(err)
			w.Write([]byte("Error: common segments in add ad delete lists.\n"))
			return
		}
	}

	// checking if segments from add list doesn't exist
	for _, segment := range request.SegmentsToAdd {
		n, err := db.DB.Count("segments", "name", segment)
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
		if n == 0 {
			logger.Logger.Log("Error: no needed segment.")
			w.Write([]byte("Error: no needed segment."))
			return
		}
	}

	// checking if segments from delete list exist
	for _, segment := range request.SegmentsToDelete {
		n, err := db.DB.Select("segments", "name", segment)
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
		if len(n) == 0 {
			logger.Logger.Log("Error: no needed segment.")
			w.Write([]byte("Error: no needed segments"))
			return
		}
	}

	// adding user in segments
	for _, val := range request.SegmentsToAdd {
		var n []string

		// checking if user is already in segment
		n, err := db.DB.Select("user_segments", "segment_name = '"+val+"' AND user_id = '"+strconv.Itoa(request.UserID)+"'", "count(*)")
		if len(n) > 0 || err != nil {
			continue
		}

		// SQL transaction to insert data in segmentHistory and userSegment
		// ------
		err = db.DB.ExecSegments([]int{request.UserID}, []string{val}, true, request.TTL)
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
		if request.TTL != "" {
			ttlUSers[request.TTL] = append(ttlUSers[request.TTL], point{request.UserID, val})
		}

		// ------
	}
	logger.Logger.Log("Segments added.")
	w.Write([]byte("Segments added.\n"))

	// deleting user from segments
	for _, val := range request.SegmentsToDelete {
		var n []string

		// checking if user is already in segment
		n, err := db.DB.Select("user_segments", "segment_name = '"+val+"' AND user_id = '"+strconv.Itoa(request.UserID)+"'", "count(*)")
		if len(n) == 0 || err != nil {
			continue
		}

		// SQL transaction to insert data in segmentHistory and userSegment
		// ------
		err = db.DB.ExecSegments([]int{request.UserID}, []string{val}, false, "")
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
		// ------
	}

	logger.Logger.Log("Segments deleted.")
	w.Write([]byte("Segments deleted.\n"))
}

// searching for active segments of user
func activeUserSegments(w http.ResponseWriter, r *http.Request) {
	var userID userRequest
	err := decode(&userID, w, r)
	if err != nil {
		return
	}
	userSegments, err := db.DB.Select("user_segments", "user_id = "+strconv.Itoa(userID.UserID), "segment_name")
	if err != nil {
		w.Write([]byte("Error: error in SQL request.\n"))
		logger.Logger.Log(err)
		return
	}
	response, err := json.Marshal(userSegments)
	if err != nil {
		w.Write([]byte("Error: error in json encoding.\n"))
		logger.Logger.Log(err)
		return
	}

	logger.Logger.Log(response)
	w.Write(response)
}

// creating file with history and response with its url
func segmentHistory(w http.ResponseWriter, r *http.Request) {
	var historyReq historyRequest
	err := decode(&historyReq, w, r)
	if err != nil {
		return
	}
	res, err := db.DB.SelectAdvanced("segment_history", historyReq.Date)
	if err != nil {
		w.Write([]byte("Error: error in SQL request.\n"))
		logger.Logger.Log(err)
		return
	}
	if len(res) == 0 {
		logger.Logger.Log("No records on " + historyReq.Date + ".")
		w.Write([]byte("No records on " + historyReq.Date + "."))
		return
	}
	f, err := os.Create("pkg/files/history/history" + historyReq.Date + ".csv")
	defer f.Close()

	if err != nil {
		w.Write([]byte("Error: error in openning file.\n"))
		logger.Logger.Log(err)
		return
	}
	for _, val := range res {
		if _, err := f.WriteString(val[0] + ";" + val[1] + ";" + val[2] + ";" + val[3] + "\n"); err != nil {
			w.Write([]byte("Error: error with file.\n"))
			logger.Logger.Log(err)
			return
		}
	}

	logger.Logger.Log("csv history created.")
	w.Write([]byte("/csvhistory?date=" + historyReq.Date))

}

// history csv response
func csvHistory(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	f, err := os.OpenFile("pkg/files/history/history"+date+".csv", os.O_RDONLY, 600)
	if err != nil {
		w.Write([]byte("Error: error with file.\n"))
		logger.Logger.Log(err)
		return
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
	w.Write(data)
}
