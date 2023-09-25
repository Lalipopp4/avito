package middleware

import (
	"net/http"

	"github.com/Lalipopp4/test_api/internal/models"
	"github.com/Lalipopp4/test_api/internal/scripts"
	"github.com/gin-gonic/gin"
)

func ValidateBookJSON(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		segment := &models.Segment{}
		if scripts.Decode(segment, c.Request.Body) != nil || segment.Name == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			// logger.Error("Bad request.")
			c.Abort()
		}
		c.Next()
	}
}
