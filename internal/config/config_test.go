package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test config structures
type SimpleConfig struct {
	Name  string `yaml:"name"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type NestedConfig struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DBConfig struct {
	Type     string       `yaml:"type"`
	SQLite   SQLiteConfig `yaml:"sqlite"`
	Postgres PGConfig     `yaml:"postgres"`
}

type SQLiteConfig struct {
	Filepath string `yaml:"filepath"`
}

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func TestLoad_SimpleConfigFromEnv(t *testing.T) {
	// Set environment variables
	t.Setenv("NAME", "test-app")
	t.Setenv("PORT", "8080")
	t.Setenv("DEBUG", "true")

	var cfg SimpleConfig
	err := Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "test-app", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
	assert.True(t, cfg.Debug)
}

func TestLoad_NestedConfigFromEnv(t *testing.T) {
	// Set nested environment variables
	t.Setenv("SERVER_HOST", "localhost")
	t.Setenv("SERVER_PORT", "3000")
	t.Setenv("DB_TYPE", "sqlite")
	t.Setenv("DB_SQLITE_FILEPATH", "/data/test.db")

	var cfg NestedConfig
	err := Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "sqlite", cfg.DB.Type)
	assert.Equal(t, "/data/test.db", cfg.DB.SQLite.Filepath)
}

func TestLoad_ConfigFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".testapp.yaml")

	configContent := `
name: "file-app"
port: 9090
debug: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Change to temp directory so config file is found
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	var cfg SimpleConfig
	err = Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "file-app", cfg.Name)
	assert.Equal(t, 9090, cfg.Port)
	assert.False(t, cfg.Debug)
}

func TestLoad_NestedConfigFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".testapp.yaml")

	configContent := `
server:
  host: "config-host"
  port: 7070
db:
  type: "postgres"
  postgres:
    host: "pg-server"
    port: 5432
    database: "testdb"
    user: "testuser"
    password: "testpass"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Change to temp directory so config file is found
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	var cfg NestedConfig
	err = Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "config-host", cfg.Server.Host)
	assert.Equal(t, 7070, cfg.Server.Port)
	assert.Equal(t, "postgres", cfg.DB.Type)
	assert.Equal(t, "pg-server", cfg.DB.Postgres.Host)
	assert.Equal(t, 5432, cfg.DB.Postgres.Port)
	assert.Equal(t, "testdb", cfg.DB.Postgres.Database)
	assert.Equal(t, "testuser", cfg.DB.Postgres.User)
	assert.Equal(t, "testpass", cfg.DB.Postgres.Password)
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".testapp.yaml")

	configContent := `
name: "file-app"
port: 9090
debug: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set environment variables to override file
	t.Setenv("NAME", "env-app")
	t.Setenv("DEBUG", "true")

	// Change to temp directory so config file is found
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	var cfg SimpleConfig
	err = Load("testapp", "", &cfg)
	require.NoError(t, err)

	// Environment should override file
	assert.Equal(t, "env-app", cfg.Name, "env should override file")
	assert.True(t, cfg.Debug, "env should override file")
	// File value should be used when no env var
	assert.Equal(t, 9090, cfg.Port, "file value when no env var")
}

func TestLoad_SpecificConfigFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-config.yaml")

	configContent := `
name: "custom-app"
port: 5555
debug: true
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var cfg SimpleConfig
	err = Load("testapp", configPath, &cfg)
	require.NoError(t, err)

	assert.Equal(t, "custom-app", cfg.Name)
	assert.Equal(t, 5555, cfg.Port)
	assert.True(t, cfg.Debug)
}

func TestLoad_MissingConfigFile(t *testing.T) {
	// No config file exists, should load from env only
	t.Setenv("NAME", "env-only")
	t.Setenv("PORT", "4000")
	t.Setenv("DEBUG", "false")

	var cfg SimpleConfig
	err := Load("nonexistent", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "env-only", cfg.Name)
	assert.Equal(t, 4000, cfg.Port)
	assert.False(t, cfg.Debug)
}

func TestLoad_InvalidConfigFile(t *testing.T) {
	// Create temporary invalid YAML file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".testapp.yaml")

	invalidYaml := `
name: "test
port: [this is invalid
`
	err := os.WriteFile(configPath, []byte(invalidYaml), 0644)
	require.NoError(t, err)

	var cfg SimpleConfig
	err = Load("testapp", configPath, &cfg)
	assert.Error(t, err)
}

type ConfigWithTags struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Ignored  string `yaml:"-"`
	Optional string `yaml:"optional,omitempty"`
}

func TestLoad_StructTags(t *testing.T) {
	// Test struct tag handling
	t.Setenv("NAME", "tagged-app")
	t.Setenv("PORT", "7777")
	t.Setenv("OPTIONAL", "present")

	var cfg ConfigWithTags
	err := Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "tagged-app", cfg.Name)
	assert.Equal(t, 7777, cfg.Port)
	assert.Equal(t, "present", cfg.Optional)
}

type ConfigWithEmbedded struct {
	ServerConfig `mapstructure:",squash"`
	Debug        bool `yaml:"debug"`
}

func TestLoad_EmbeddedStruct(t *testing.T) {
	// Test embedded/squashed struct
	t.Setenv("HOST", "embedded-host")
	t.Setenv("PORT", "6666")
	t.Setenv("DEBUG", "true")

	var cfg ConfigWithEmbedded
	err := Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "embedded-host", cfg.Host)
	assert.Equal(t, 6666, cfg.Port)
	assert.True(t, cfg.Debug)
}

type ComplexConfig struct {
	Level1 Level1Config `yaml:"level1"`
}

type Level1Config struct {
	Level2 Level2Config `yaml:"level2"`
	Name   string       `yaml:"name"`
}

type Level2Config struct {
	Value int `yaml:"value"`
}

func TestLoad_DeepNesting(t *testing.T) {
	// Test nested configuration (2 levels deep)
	t.Setenv("LEVEL1_NAME", "first")
	t.Setenv("LEVEL1_LEVEL2_VALUE", "42")

	var cfg ComplexConfig
	err := Load("testapp", "", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "first", cfg.Level1.Name)
	assert.Equal(t, 42, cfg.Level1.Level2.Value)
}

func TestLoad_MixedSources(t *testing.T) {
	// Create config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".testapp.yaml")

	configContent := `
server:
  host: "file-host"
  port: 8000
db:
  type: "file-db"
  sqlite:
    filepath: "/file/path.db"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Override some values with environment variables
	t.Setenv("SERVER_PORT", "9000")
	t.Setenv("DB_SQLITE_FILEPATH", "/env/path.db")

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	var cfg NestedConfig
	err = Load("testapp", "", &cfg)
	require.NoError(t, err)

	// File value (no env override)
	assert.Equal(t, "file-host", cfg.Server.Host)
	// Env override
	assert.Equal(t, 9000, cfg.Server.Port)
	assert.Equal(t, "/env/path.db", cfg.DB.SQLite.Filepath)
	// File value (no env override)
	assert.Equal(t, "file-db", cfg.DB.Type)
}
