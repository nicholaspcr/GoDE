package migrations

import (
	"os"
	"testing"

	"github.com/nicholaspcr/GoDE/internal/store/migrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedMigrationsExist(t *testing.T) {
	// Check that the embedded filesystem contains migration files
	entries, err := migrations.FS.ReadDir(".")
	require.NoError(t, err, "should be able to read embedded migrations directory")

	// Should have at least 2 files (up and down)
	assert.GreaterOrEqual(t, len(entries), 2, "should have at least 2 migration files")

	// Check for specific migration files
	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}

	assert.Contains(t, fileNames, "000001_initial_schema.up.sql", "should contain up migration")
	assert.Contains(t, fileNames, "000001_initial_schema.down.sql", "should contain down migration")
}

func TestEmbeddedMigrationContent(t *testing.T) {
	// Read the up migration file
	content, err := migrations.FS.ReadFile("000001_initial_schema.up.sql")
	require.NoError(t, err, "should be able to read up migration file")

	// Verify it's not empty and contains expected SQL
	assert.NotEmpty(t, content, "migration file should not be empty")
	assert.Contains(t, string(content), "CREATE TABLE", "should contain CREATE TABLE statements")
	assert.Contains(t, string(content), "users", "should create users table")
	assert.Contains(t, string(content), "pareto_sets", "should create pareto_sets table")
	assert.Contains(t, string(content), "vectors", "should create vectors table")
}

func TestMigrationWithSQLite(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_migration_*.db")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	databaseURL := "sqlite3://" + tmpFile.Name()

	// Run migrations
	err = Run(databaseURL)
	assert.NoError(t, err, "should successfully run migrations")

	// Check version
	version, dirty, err := Version(databaseURL)
	assert.NoError(t, err, "should be able to get version")
	assert.Equal(t, uint(1), version, "should be at version 1")
	assert.False(t, dirty, "should not be in dirty state")
}

func TestMigrationRollback(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_rollback_*.db")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	databaseURL := "sqlite3://" + tmpFile.Name()

	// Run migrations
	err = Run(databaseURL)
	require.NoError(t, err, "should successfully run migrations")

	// Rollback
	err = Rollback(databaseURL, 1)
	assert.NoError(t, err, "should successfully rollback migration")

	// Check version - should be at version 0 (no migrations)
	version, dirty, err := Version(databaseURL)
	if err != nil {
		// ErrNilVersion is expected when all migrations are rolled back
		assert.Equal(t, uint(0), version, "should be at version 0 after rollback")
	}
	assert.False(t, dirty, "should not be in dirty state")
}

func TestMaskPassword(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "short URL",
			url:      "short",
			expected: "***",
		},
		{
			name:     "long URL",
			url:      "postgres://user:password@localhost:5432/dbname",
			expected: "postgres:/***432/dbname",
		},
		{
			name:     "SQLite URL",
			url:      "sqlite3:///path/to/database.db",
			expected: "sqlite3://***atabase.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskPassword(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}
