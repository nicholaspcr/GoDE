package routes

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	dataMap   map[string]any
	dataArray []any
)

func leftNavbar() string {
	var buff bytes.Buffer
	t, ok := templates["navbar-set.html"]
	if !ok {
		slog.Error("Error fetching template: navbar-set")
		return ""
	}
	err := t.Execute(&buff, dataArray{
		dataMap{
			"URL":           "/home",
			"Name":          "Home",
			"ShouldPushURL": true,
			"PushURL":       "/",
		},
	})
	if err != nil {
		slog.Error("Error parsing template: ", err)
		return ""
	}
	return buff.String()
}

func rightNavbar() string {
	var buff bytes.Buffer
	t, ok := templates["navbar-set.html"]
	if !ok {
		slog.Error("Error fetching template: navbar-set")
		return ""
	}
	err := t.Execute(&buff, dataArray{
		dataMap{
			"URL":  "/login",
			"Name": "Login",
		},
		dataMap{
			"URL":  "/register",
			"Name": "Register",
		},
	})
	if err != nil {
		slog.Error("Error parsing template: ", err)
		return ""
	}
	return buff.String()
}

// baseRoutes defines the base routes for the API.
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
