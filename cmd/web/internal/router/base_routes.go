package router

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	dataMap   map[string]any
	dataArray []any
)

// baseRoutes - base routes for the API.
func baseRoutes(r *echo.Group) {
	r.GET("/", func(c echo.Context) error {
		t, ok := templates["index.html"]
		if !ok {
			return c.String(http.StatusNoContent, "Template not found")
		}
		data := dataMap{
			"LeftNav":  dataArray{template.HTML(leftNavbar())},
			"RightNav": dataArray{template.HTML(rightNavbar())},
		}
		if err := t.Execute(c.Response().Writer, data); err != nil {
			c.Response().Status = http.StatusInternalServerError
		}
		c.Response().Status = http.StatusOK
		return nil
	})

	r.GET("/home", func(c echo.Context) error {
		t, ok := templates["home.html"]
		if !ok {
			return c.String(http.StatusNoContent, "Template not found")
		}
		if err := t.Execute(c.Response().Writer, nil); err != nil {
			c.Response().Status = http.StatusInternalServerError
		}
		c.Response().Status = http.StatusOK
		return nil
	})
}
