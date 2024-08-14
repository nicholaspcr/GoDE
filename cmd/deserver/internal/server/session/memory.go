package session

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type inMemorySessionStore struct {
	*expirable.LRU[string, struct{}]
}

// NewInMemoryStore returns a in memory implementation of the Store operations.
func NewInMemoryStore() Store {
	return &inMemorySessionStore{
		LRU: expirable.NewLRU[string, struct{}](0, nil, 12*time.Hour),
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
