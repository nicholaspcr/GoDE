package gorm

// Config of the gorm implementation of the store.Interfaces.
type Config struct {
	UseMemory bool   `json:"use-memory" yaml:"use-memory"`
	DNS       string `json:"dns" yaml:"dns"`
}
