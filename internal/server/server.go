// Package server includes the implementation of the API services.
package server

import "github.com/nicholaspcr/GoDE/pkg/api"

type server struct {
	*userServices
}

// Server is the interface that provides the API services.
type Server interface {
	api.UserBaseServicesServer
}

// NewServer creates a new server.
func NewServer() Server {
	return &server{
		userServices: &userServices{},
	}
}
