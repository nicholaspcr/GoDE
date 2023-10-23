package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// baseRoutes defines the base routes for the API.
func baseRoutes(r *echo.Group) {
	r.GET("/", func(c echo.Context) error {
		t, ok := templates["index.html"]
		if !ok {
			return c.String(http.StatusNoContent, "Template not found")
		}
		if err := t.Execute(c.Response().Writer, nil); err != nil {
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

	r.GET("/login", func(c echo.Context) error {
		t, ok := templates["login.html"]
		if !ok {
			return c.String(http.StatusNoContent, "Template not found")
		}
		if err := t.Execute(c.Response().Writer, nil); err != nil {
			c.Response().Status = http.StatusInternalServerError
		}
		c.Response().Status = http.StatusOK
		return nil
	})

	r.GET("/register", func(c echo.Context) error {
		t, ok := templates["register.html"]
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
