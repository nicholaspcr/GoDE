package server

import "github.com/nicholaspcr/GoDE/internal/store"

type serverOpts func(*server)

// WithStore sets the store for the server instance.
func WithStore(st store.Store) serverOpts {
	return func(s *server) { s.st = st }
}

// WithConfig parser
func WithConfig(cfg Config) serverOpts {
	return func(s *server) { s.cfg = cfg }
}
