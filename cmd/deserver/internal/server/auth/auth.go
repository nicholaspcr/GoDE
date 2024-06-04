// Package auth contains the authentication methods and middlewares.
package auth

import (
	"encoding/base64"
	"net/http"
	"net/mail"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// Adds the /register route to the HTTP server
func RegisterHandler(mux *runtime.ServeMux, st store.UserOperations) error {
	return mux.HandlePath(
		"POST",
		"/register",
		func(
			w http.ResponseWriter,
			r *http.Request,
			pathParams map[string]string,
		) {
			ctx := r.Context()
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			name := r.Form.Get("name")
			username := r.Form.Get("username")
			password := r.Form.Get("password")

			addr, err := mail.ParseAddress(username)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := st.CreateUser(ctx, &api.User{
				Ids:      &api.UserIDs{Email: addr.Address},
				Name:     name,
				Password: password,
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})
}

// Adds the /login route to the HTTP server
func LoginHandler(mux *runtime.ServeMux, st store.UserOperations) error {
	return mux.HandlePath(
		"POST",
		"/login",
		func(
			w http.ResponseWriter,
			r *http.Request,
			pathParams map[string]string,
		) {
			ctx := r.Context()
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			username := r.Form.Get("username")
			password := r.Form.Get("password")

			addr, err := mail.ParseAddress(username)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			usr, err := st.GetUser(ctx, &api.UserIDs{Email: addr.Address})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if usr.Password != password {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Add authorization token in the response header.
			authValue := username + ":" + password
			authToken := base64.StdEncoding.EncodeToString([]byte(authValue))
			w.Header().Add("Authorization", "Basic "+authToken)

			w.WriteHeader(http.StatusOK)
		})
}
