package store

// Store contains the methods to interact with the database
type Store interface {
	User
}

type store struct {
	*user
}

// New returns a new Store instance
func New() Store {
	return &store{
		user: &user{},
	}
}
