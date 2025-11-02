package session

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

const (
	// DefaultMaxSessions is the maximum number of concurrent sessions allowed.
	// This prevents memory exhaustion from session storage.
	DefaultMaxSessions = 10000

	// DefaultSessionTTL is the default time-to-live for sessions.
	// Sessions are automatically cleaned up after this duration.
	DefaultSessionTTL = 12 * time.Hour
)

type inMemorySessionStore struct {
	*expirable.LRU[string, struct{}]
}

// NewInMemoryStore returns a in memory implementation of the Store operations.
// It uses an LRU cache with automatic expiration to prevent memory leaks.
// - maxSize: maximum number of concurrent sessions (use DefaultMaxSessions if unsure)
// - ttl: how long sessions remain valid (use DefaultSessionTTL if unsure)
func NewInMemoryStore(maxSize int, ttl time.Duration) Store {
	if maxSize <= 0 {
		maxSize = DefaultMaxSessions
	}
	if ttl <= 0 {
		ttl = DefaultSessionTTL
	}

	return &inMemorySessionStore{
		LRU: expirable.NewLRU[string, struct{}](maxSize, nil, ttl),
	}
}

func (st *inMemorySessionStore) Add(k string) {
	st.LRU.Add(k, struct{}{})
}

func (st *inMemorySessionStore) Get(k string) bool {
	_, ok := st.LRU.Get(k)
	return ok
}

func (st *inMemorySessionStore) Remove(k string) bool {
	return st.LRU.Remove(k)
}
