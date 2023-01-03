package middleware

import (
	"message/internal/app"

	"github.com/gin-gonic/gin"
)

func Cors(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "51.250.66.160")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Auth-Token, Content-Type, Content-Length, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
	}
}
