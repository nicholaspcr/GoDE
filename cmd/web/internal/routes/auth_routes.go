package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

// accountRoutes - auth routes for the API.
func accountRoutes(r *echo.Group) {
	r.GET("/login", func(c echo.Context) error {
		slog.Info("NICK - response header",
			slog.String("header", c.Response().Header().Get("go_web_path_validation_header")),
		)
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
