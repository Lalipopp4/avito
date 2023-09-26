package rest

import (
	"context"
	"net/http"

	"github.com/Lalipopp4/test_api/internal/models"
	"github.com/Lalipopp4/test_api/internal/scripts"
	"github.com/Lalipopp4/test_api/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

func handle(rAPI *restAPI) *gin.Engine {

	rtr := gin.Default()
	rtr.Use(
		gin.Logger(),
		gin.Recovery(),
	)
	segments := rtr.Group("segment")
	segments.Use(middleware.ValidateSegment(rAPI.logger))
	{
		segments.POST("/segment/add", rAPI.addSegment)
		segments.POST("/segment/delete", rAPI.deleteSegment)

	}

	rtr.POST("/segment/history", middleware.ValidateDate(rAPI.logger), rAPI.getHistoryByDate)
	rtr.GET("/segment/history", rAPI.getCSV)

	users := rtr.Group("user")
	users.Use(middleware.ValidateUser(rAPI.logger))
	{
		rtr.GET("/user/segments", rAPI.getSegmentsByUser)
		rtr.POST("/user/add", rAPI.addUser)
	}

	return rtr
}

func (rAPI *restAPI) Run() error {
	go rAPI.segmentService.CheckTTL()
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
		rAPI.logger.Error(err.Error())
		return
	}

	err = rAPI.segmentService.AddSegment(context.Background(), segment)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error(err.Error())
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
		rAPI.logger.Error(err.Error())
		return
	}

	err = rAPI.segmentService.DeleteSegment(context.Background(), segment)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error(err.Error())
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
		rAPI.logger.Error(err.Error())
		return
	}

	err = rAPI.userService.AddUser(context.Background(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error(err.Error())
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
		rAPI.logger.Error(err.Error())
		return
	}

	segments, err := rAPI.userService.GetSegmentsByUser(context.Background(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error(err.Error())
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
		rAPI.logger.Error(err.Error())
		return
	}

	link, err := rAPI.segmentService.GetHistoryByDate(context.Background(), date)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		rAPI.logger.Error(err.Error())
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

	rAPI.logger.Info("CSV file on " + date + " created and sent.")
	c.JSON(http.StatusCreated, data)
}
