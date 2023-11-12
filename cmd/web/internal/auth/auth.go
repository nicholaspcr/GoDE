// Package auth provides the necessary methods for authentication an user on the
// web server.
package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

const (
	// TODO: Make it configurable via env variables.
	jwtCookieName     = "token"
	jwtExpireDuration = 72 * time.Hour
	jwtSecret         = "secret"
)

// JWTClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JWTClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

// FromClaims returns a JWTClaims struct from a jwt.MapClaims.
func GetClaims(c echo.Context) JWTClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return JWTClaims{
		Name:  claims["name"].(string),
		Admin: claims["admin"].(bool),
	}
}

// Login validates if the user exists within the database and returns the claims
// based on the user properties.
func Login(c echo.Context) (*http.Cookie, error) {
	if err := c.Request().ParseForm(); err != nil {
		return nil, err
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	// Throws unauthorized error
	if email != "john@gmail.com" || password != "12345" {
		slog.With(
			slog.String("username", email),
			slog.String("password", password),
		).Error("Invalid username or password")
		return nil, echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JWTClaims{
		"Jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpireDuration)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token.
	value, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:    jwtCookieName,
		Value:   value,
		Expires: time.Now().Add(jwtExpireDuration),
	}, nil
}
