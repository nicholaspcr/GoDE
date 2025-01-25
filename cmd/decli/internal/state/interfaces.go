package state

// Operations related to the State of the CLI application.g
type Operations interface {
	AuthTokenOperations
}

// AuthTokenOperations is the set of operations related to the Authentication
// token received by the login command in the CLI.
type AuthTokenOperations interface {
	// GetAuthToken gets the latest authentication token from the state store.
	GetAuthToken() (string, error)
	// InvalidateAuthToken invalidates the latest authentication token.
	InvalidateAuthToken() error
	// SaveAuthToken saves the authentication token as the latest token.
	SaveAuthToken(string) error
}
