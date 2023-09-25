package rest

import "github.com/gin-gonic/gin"

type Server interface {
	Run() error
	addSegment(c *gin.Context)
	deleteSegment(c *gin.Context)
	getSegmentsByUser(c *gin.Context)
	addUser(c *gin.Context)
	getHistoryByDate(c *gin.Context)
	getCSV(c *gin.Context)
	Stop() error
}
