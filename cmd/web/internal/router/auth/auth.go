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
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
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
	Email string `json:"email"`
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

	conn, err := grpc.Dial(DB_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to connect to deserver",
			slog.String("err_msg", err.Error()),
			slog.String("db_url", DB_URL),
		)
		return nil, err
	}
	usrClient := api.NewUserServiceClient(conn)
	res, err := usrClient.Get(
		c.Request().Context(),
		&api.UserServiceGetRequest{UserIds: &api.UserIDs{Email: email}},
	)
	if err != nil {
		return nil, err
	}

	// TODO: Remove this.
	slog.With(slog.String("usr", res.User.String())).Debug("User found")

	if res.User.Password != password {
		slog.With("usr_password", res.User.Password, "password", password).Debug("User password does not match.")
		return nil, echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JWTClaims{
		res.User.Ids.Email,
		res.User.Name,
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

func Register(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}

	email := c.FormValue("email")
	password := c.FormValue("password")
	name := c.FormValue("name")

	newUser := &api.User{
		Ids:      &api.UserIDs{Email: email},
		Name:     name,
		Password: password,
	}

	conn, err := grpc.Dial(DB_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.With(
			slog.String("err_msg", err.Error()),
			slog.String("db_url", DB_URL),
		).Error("failed to connect to deserver")
		return err
	}

	usrClient := api.NewUserServiceClient(conn)
	_, err = usrClient.Create(
		c.Request().Context(),
		&api.UserServiceCreateRequest{User: newUser},
	)
	if err != nil {
		return err
	}

	return nil
}

func ForgotPassword(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}

	email := c.FormValue("email")

	usrIDs := &api.UserIDs{Email: email}
	_ = usrIDs

	//conn, err := grpc.Dial(DB_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	slog.With(
	//		slog.String("err_msg", err.Error()),
	//		slog.String("db_url", DB_URL),
	//	).Error("failed to connect to deserver")
	//	return err
	//}

	//usrClient := api.NewUserServiceClient(conn)
	//_, err = usrClient.Create(
	//	c.Request().Context(),
	//	&api.UserServiceCreateRequest{User: newUser},
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}
