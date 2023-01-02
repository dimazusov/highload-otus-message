package http

import (
	"io/ioutil"
	"log"
	"message/internal/app"
	"message/internal/server/http/handlers/user"
	"message/internal/server/http/middleware"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// @title Swagger API
// @version 1.0
// @description message api
func NewGinRouter(app *app.App) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/*.html")

	frontendGroup := router.Group("/")
	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{}) })
	router.GET("/auth", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{}) })
	router.GET("/registration", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{}) })
	router.GET("/manifest.json", func(c *gin.Context) {
		file, err := os.Open("web/manifest.json")
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		data := strings.ReplaceAll(string(b), "\n", "")
		data = strings.ReplaceAll(data, "\\\"", "\"")
		c.Data(http.StatusOK, "application/json", []byte(data))
	})
	frontendGroup.StaticFS("/static", http.Dir("web/static"))

	authGroup := router.Group("/api/v1").Use(middleware.Cors(app))

	authGroup.OPTIONS("/auth", func(c *gin.Context) { c.AbortWithStatus(http.StatusOK) })
	authGroup.OPTIONS("/registration", func(c *gin.Context) { c.AbortWithStatus(http.StatusOK) })

	authGroup.POST("/auth", func(c *gin.Context) { user.AuthHandler(c, app) })
	authGroup.POST("/registration", func(c *gin.Context) { user.RegisterHandler(c, app) })

	apiGroup := router.Group("/api/v1").
		Use(middleware.Cors(app))
	//Use(middleware.Auth(app))

	authGroup.OPTIONS("/user/:id", func(c *gin.Context) { c.AbortWithStatus(http.StatusOK) })

	apiGroup.GET("/users", func(c *gin.Context) { user.GetUsersHandler(c, app) })
	apiGroup.GET("/user/:id", func(c *gin.Context) { user.GetUserHandler(c, app) })
	apiGroup.PUT("/user", func(c *gin.Context) { user.UpdateUserHandler(c, app) })
	apiGroup.POST("/user", func(c *gin.Context) { user.CreateUserHandler(c, app) })
	apiGroup.DELETE("/user/:id", func(c *gin.Context) { user.DeleteUserHandler(c, app) })

	return router
}
