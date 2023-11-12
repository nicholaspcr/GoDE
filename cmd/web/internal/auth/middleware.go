package auth

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Middleware returns a middleware function that validates JWT token.
func Middleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		TokenLookup: "cookie:" + jwtCookieName,
		SigningKey:  []byte(jwtSecret),
	})
}
