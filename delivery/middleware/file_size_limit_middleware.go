package middleware

import (
	"instructor-led-app/shared/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _ := c.FormFile("image")

		if file != nil && file.Size > maxSize {
			common.SendErrorResponse(c, http.StatusBadRequest, "File size exceeds the limit")
			c.Abort()
			return
		}

		c.Next()
	}
}
