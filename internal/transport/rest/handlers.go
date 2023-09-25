package rest

import (
	"context"
	"net/http"

	"github.com/Lalipopp4/test_api/internal/models"
	"github.com/Lalipopp4/test_api/internal/scripts"
	"github.com/gin-gonic/gin"
)

func handle(rAPI Server) *gin.Engine {

	rtr := gin.Default()
	rtr.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	rtr.POST("/segment/add", rAPI.addSegment)

	rtr.POST("/segment/delete", rAPI.deleteSegment)

	rtr.POST("/user/segments", rAPI.getSegmentsByUser)

	rtr.POST("/user/add", rAPI.addUser)

	rtr.POST("/segment/history", rAPI.getHistoryByDate)

	rtr.GET("/segment/history", rAPI.getCSV)

	return rtr
}

func (rAPI *restAPI) Run() error {
	return rAPI.httpServer.ListenAndServe()
}

func (rAPI *restAPI) Stop() error {
	return nil
}

func (rAPI *restAPI) addSegment(c *gin.Context) {
	segment := &models.Segment{}
	err := scripts.Decode(&segment, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}

	err = rAPI.segmentService.AddSegment(context.Background(), segment)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	rAPI.logger.Info(segment.Name + " added successfully.")
	c.JSON(http.StatusCreated, segment.Name+" added successfully.")
}

func (rAPI *restAPI) deleteSegment(c *gin.Context) {
	segment := &models.Segment{}
	err := scripts.Decode(&segment, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}

	err = rAPI.segmentService.DeleteSegment(context.Background(), segment)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	rAPI.logger.Info(segment.Name + " deleted successfully.")
	c.JSON(http.StatusCreated, segment.Name+" deleted successfully.")
}

func (rAPI *restAPI) addUser(c *gin.Context) {
	user := &models.UserRequest{}
	err := scripts.Decode(&user, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	err = rAPI.userService.AddUser(context.Background(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	rAPI.logger.Info("Segments were added and deleted for user successfully.")
	c.JSON(http.StatusCreated, "Segments were added and deleted for user successfully.")
}

func (rAPI *restAPI) getSegmentsByUser(c *gin.Context) {
	user := &models.User{}
	err := scripts.Decode(&user, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	segments, err := rAPI.userService.GetSegmentsByUser(context.Background(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	rAPI.logger.Info("Active user segments sent.")
	c.JSON(http.StatusCreated, segments)
}

func (rAPI *restAPI) getHistoryByDate(c *gin.Context) {
	date := &models.Date{}
	err := scripts.Decode(&date, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	link, err := rAPI.segmentService.GetHistoryByDate(context.Background(), date)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	rAPI.logger.Info("History file on " + date.Date + " created and sent.")
	c.JSON(http.StatusCreated, link)
}

func (rAPI *restAPI) getCSV(c *gin.Context) {
	date := c.Params.ByName("date")
	data, err := rAPI.segmentService.GetCSV(context.Background(), date)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error("Error in request.")
		return
	}
	c.JSON(http.StatusCreated, data)
}
