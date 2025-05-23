// Package session defines and implements the operations necessary for a session
// store, this is used for authentication of users.
package session

// Store handles the operations to handles multiples users sessions.
type Store interface {
	Add(k string)
	Get(s string) bool
	Remove(s string) bool
}
