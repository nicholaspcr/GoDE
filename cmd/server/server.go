package main

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/server"
)

func main() {
	ctx := context.TODO()

	srv, err := server.New(ctx, server.DefaultConfig())
	if err != nil {
		panic(err)
	}

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}
}
