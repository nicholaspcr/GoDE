package routes

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/labstack/echo/v4"
)

const templatesDir = "templates"

var (
	//go:embed templates/*
	files     embed.FS
	templates map[string]*template.Template
)

// Route defines a set of routes for the API. It attaches the routes to the
// echo.RouterGroup provided.
type Route func(*echo.Group)

// RouteGroups is a list of predefined routes for the API.
var RouteGroups = map[string]Route{
	"":     baseRoutes,
	"auth": authRoutes,
}

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
}
