package routes

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed templates
var templateFS embed.FS

// Route defines a set of routes for the API. It attaches the routes to the gin.RouterGroup provided.
type Route func(*gin.RouterGroup)

// Routes is a list of predefined routes for the API.
var Routes = []Route{
	baseRoutes,
}

// baseRoutes defines the base routes for the API.
func baseRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS, "templates/index.html"))
		t.Execute(c.Writer, nil)
	})

	r.GET("/name", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS, "templates/name.html"))
		t.Execute(c.Writer, gin.H{"name": "Gopher"})
	})
}
