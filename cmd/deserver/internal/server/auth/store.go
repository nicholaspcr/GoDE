package auth

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type SessionStore interface {
	Add(k string)
	Get(s string) bool
	Remove(s string) bool
}

type inMemorySessionStore struct {
	*expirable.LRU[string, struct{}]
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

func NewInMemorySessionStore() SessionStore {
	return &inMemorySessionStore{
		LRU: expirable.NewLRU[string, struct{}](0, nil, time.Minute),
	}
}
