// Package handlers contains the HTTP handlers for the API.
package handlers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/store"
	"google.golang.org/grpc"
)

// Handler creates a wrapper for each individualized wrapper.
type Handler interface {
	SetStore(store.Store)
	RegisterService(*grpc.Server)
	RegisterHTTPHandler(
		context.Context, *runtime.ServeMux, string, []grpc.DialOption,
	) error
}
