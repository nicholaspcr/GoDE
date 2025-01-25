// Package state manages the state of the CLI.
package state

// State of the CLI application.
type State struct {
	// AuthToken is the token received in the authentication operation.
	AuthToken string
}
