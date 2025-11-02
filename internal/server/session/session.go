// Package session defines and implements the operations necessary for a session
// store, this is used for authentication of users.
//
// NOTE: This package is currently unused as the system uses stateless JWT authentication.
// Session storage is provided for future extensibility if needed (e.g., refresh tokens,
// session-based features, etc.). The current JWT-only approach is preferred for its
// stateless nature and better horizontal scalability.
package session

// Store handles the operations to handles multiples users sessions.
type Store interface {
	Add(k string)
	Get(s string) bool
	Remove(s string) bool
}
