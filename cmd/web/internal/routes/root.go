package routes

import "github.com/gin-gonic/gin"

// Route defines a set of routes for the API. It attaches the routes to the gin.RouterGroup provided.
type Route func(*gin.RouterGroup)

// Routes is a list of predefined routes for the API.
var Routes = []Route{
	BaseRoutes,
}

// BaseRoutes defines the base routes for the API.
func BaseRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
