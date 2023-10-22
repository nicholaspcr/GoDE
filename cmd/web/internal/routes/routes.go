package routes

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/gin-gonic/gin"
)

const templatesDir = "templates"

var (
	//go:embed templates/*
	files     embed.FS
	templates map[string]*template.Template
)

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		panic(err)
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(
			files,
			templatesDir+"/"+tmpl.Name(),
		)
		if err != nil {
			panic(err)
		}

		templates[tmpl.Name()] = pt
	}

	srvPaths := []string{}
	for k := range templates {
		srvPaths = append(srvPaths, k)
	}
	fmt.Println("Serving templates:", srvPaths)
}

// Route defines a set of routes for the API. It attaches the routes to the
// gin.RouterGroup provided.
type Route func(*gin.RouterGroup)

// Routes is a list of predefined routes for the API.
var Routes = []Route{
	baseRoutes,
}

// baseRoutes defines the base routes for the API.
func baseRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		t := templates["index.html"]
		if err := t.Execute(c.Writer, nil); err != nil {
			c.AbortWithError(500, err)
			return
		}
	})

	r.GET("/login", func(c *gin.Context) {
		t := templates["login.html"]
		if err := t.Execute(c.Writer, nil); err != nil {
			c.AbortWithError(500, err)
			return
		}
	})

	r.GET("/register", func(c *gin.Context) {
		t := templates["register.html"]
		if err := t.Execute(c.Writer, nil); err != nil {
			c.AbortWithError(500, err)
			return
		}
	})

	r.GET("/home", func(c *gin.Context) {
		t := templates["home.html"]
		if err := t.Execute(c.Writer, nil); err != nil {
			c.AbortWithError(500, err)
			return
		}
	})
}
