package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicholaspcr/GoDE/cmd/web/internal/router/auth"
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

	r.POST("/register", func(c echo.Context) error {
		if err := auth.Register(c); err != nil {
			return err
		}
		c.Response().Status = http.StatusOK
		return nil
	})

	r.GET("/forgot-password", func(c echo.Context) error {
		t, ok := templates["forgot-password.html"]
		if !ok {
			return c.String(http.StatusNoContent, "Template not found")
		}
		if err := t.Execute(c.Response().Writer, nil); err != nil {
			c.Response().Status = http.StatusInternalServerError
		}

		c.Response().Status = http.StatusOK
		return nil
	})

	r.POST("/forgot-password", func(c echo.Context) error {
		if err := auth.ForgotPassword(c); err != nil {
			return err
		}
		c.Response().Header().Add("HX-Redirect", "/")
		c.Response().Status = http.StatusPermanentRedirect
		return nil
	})
}
