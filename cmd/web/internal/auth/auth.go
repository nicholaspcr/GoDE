// Package auth provides the necessary methods for authentication an user on the
// web server.
package auth

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/nicholaspcr/GoDE/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// TODO: Make it configurable via env variables.
	jwtCookieName     = "token"
	jwtExpireDuration = 72 * time.Hour
	jwtSecret         = "secret"
)

// DB_URL is the database url.
var DB_URL = os.Getenv("DB_URL")

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

	conn, err := grpc.Dial(
		DB_URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	usrClient := api.NewUserServicesClient(conn)
	usr, err := usrClient.Get(c.Request().Context(), &api.UserIDs{Email: email})
	if err != nil {
		return nil, err
	}

	// TODO: Remove this.
	slog.Debug("User found: %v", usr)
	if usr.Password != password {
		slog.With(
			"usr_password", usr.Password,
			"password", password,
		).Debug("User password does not match.")
		return nil, echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JWTClaims{
		usr.Ids.Email,
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
