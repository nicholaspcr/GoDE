package sqlite

// Config contains SQLite storage configuration for the CLI client state.
type Config struct {
	// Provider selects either 'memory' or 'file' as the storage location.
	Provider string
	// Filepath of file in which data will be saved.
	Filepath string
}
