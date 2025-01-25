package sqlite

type Config struct {
	// Provider selects either 'memory' or 'file' as the storage location.
	Provider string
	// Filepath of file in which data will be saved.
	Filepath string
}
