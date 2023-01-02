package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"message/internal/app"
	"message/internal/domain/auth_token"
)

func Auth(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		values, ok := c.Request.Header["X-Auth-Token"]
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if len(values) == 0 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		_, err := app.Domain.AuthToken.Service.Parse(context.Background(), values[0])
		if err != nil {
			if errors.Is(err, auth_token.ErrWrongToken) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Next()
	}
}

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
