package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicholaspcr/GoDE/cmd/web/internal/auth"
)

// authRoutes - auth routes for the API.
func authRoutes(r *echo.Group) {
	r.GET("/test", func(c echo.Context) error {
		claims := auth.GetClaims(c)
		return c.String(
			http.StatusOK,
			fmt.Sprintf(
				"Welcome %s, admin state is: %t",
				claims.Name, claims.Admin,
			),
		)
	}, auth.Middleware())

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

	r.POST("/login", func(c echo.Context) error {
		cookie, err := auth.Login(c)
		if err != nil {
			return err
		}

		c.SetCookie(cookie)
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
