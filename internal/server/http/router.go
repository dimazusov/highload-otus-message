package http

import (
	"message/internal/app"
	"message/internal/server/http/handlers/message"
	"message/internal/server/http/middleware"

	"github.com/gin-gonic/gin"
)

// @title Swagger API
// @version 1.0
// @description message api
func NewGinRouter(app *app.App) *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api/v1").Use(middleware.Cors(app))

	apiGroup.GET("/messages", func(c *gin.Context) { message.GetMessagesHandler(c, app) })
	apiGroup.GET("/message/:id", func(c *gin.Context) { message.GetMessageHandler(c, app) })
	apiGroup.PUT("/message", func(c *gin.Context) { message.UpdateMessageHandler(c, app) })
	apiGroup.POST("/message/send", func(c *gin.Context) { message.SendMessageHandler(c, app) })
	apiGroup.DELETE("/message/:id", func(c *gin.Context) { message.DeleteMessageHandler(c, app) })

	return router
}
