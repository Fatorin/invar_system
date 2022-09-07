package middlewares

import (
	"invar/status"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	json     = "application/json"
	formData = "multipart/form-data"
)

func RequestSizeLimit() func(c *gin.Context) {
	return func(c *gin.Context) {
		var maxBytes int64 = 1024 * 1024 * 100 // 100MB
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, maxBytes)
		contentType := c.ContentType()
		switch contentType {
		case json:
			c.Next()
		case formData:
			err := c.Request.ParseMultipartForm(1024)
			if err != nil {
				logrus.Error(err)
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					status.RespStatus: status.NewResponse(status.TooLarge),
				})
				c.Abort()
			}
		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				status.RespStatus: status.NewResponse(status.BadRequest),
			})
			c.Abort()
		}

		c.Next()
	}
}
