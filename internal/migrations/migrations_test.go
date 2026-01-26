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

	// Should have at least 2 files (up and down) for each of 7 migrations + embed.go
	assert.GreaterOrEqual(t, len(entries), 14, "should have at least 14 migration files (7 up + 7 down)")

	// Check for specific migration files
	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}

	// Verify all 7 migrations exist (up and down)
	expectedMigrations := []string{
		"000001_initial_schema.up.sql",
		"000001_initial_schema.down.sql",
		"000002_add_executions_and_indices.up.sql",
		"000002_add_executions_and_indices.down.sql",
		"000003_add_user_created_index.up.sql",
		"000003_add_user_created_index.down.sql",
		"000004_add_deleted_at_columns.up.sql",
		"000004_add_deleted_at_columns.down.sql",
		"000005_add_max_objs_to_pareto.up.sql",
		"000005_add_max_objs_to_pareto.down.sql",
		"000006_add_updated_at_to_vectors.up.sql",
		"000006_add_updated_at_to_vectors.down.sql",
		"000007_add_execution_metadata.up.sql",
		"000007_add_execution_metadata.down.sql",
	}

	for _, expected := range expectedMigrations {
		assert.Contains(t, fileNames, expected, "should contain migration file: %s", expected)
	}
}

func TestEmbeddedMigrationContent(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		contains []string
	}{
		{
			name: "000001_initial_schema.up.sql",
			file: "000001_initial_schema.up.sql",
			contains: []string{
				"CREATE TABLE",
				"users",
				"pareto_sets",
				"vectors",
			},
		},
		{
			name: "000002_add_executions_and_indices.up.sql",
			file: "000002_add_executions_and_indices.up.sql",
			contains: []string{
				"CREATE TABLE",
				"executions",
				"CREATE INDEX",
			},
		},
		{
			name: "000003_add_user_created_index.up.sql",
			file: "000003_add_user_created_index.up.sql",
			contains: []string{
				"CREATE INDEX",
				"executions",
				"pareto_sets",
				"created_at",
			},
		},
		{
			name: "000004_add_deleted_at_columns.up.sql",
			file: "000004_add_deleted_at_columns.up.sql",
			contains: []string{
				"ALTER TABLE",
				"deleted_at",
			},
		},
		{
			name: "000005_add_max_objs_to_pareto.up.sql",
			file: "000005_add_max_objs_to_pareto.up.sql",
			contains: []string{
				"ALTER TABLE",
				"pareto_sets",
				"max_objs_json",
			},
		},
		{
			name: "000006_add_updated_at_to_vectors.up.sql",
			file: "000006_add_updated_at_to_vectors.up.sql",
			contains: []string{
				"ALTER TABLE",
				"vectors",
				"updated_at",
			},
		},
		{
			name: "000007_add_execution_metadata.up.sql",
			file: "000007_add_execution_metadata.up.sql",
			contains: []string{
				"ALTER TABLE",
				"executions",
				"algorithm",
				"variant",
				"problem",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := migrations.FS.ReadFile(tt.file)
			require.NoError(t, err, "should be able to read migration file")

			assert.NotEmpty(t, content, "migration file should not be empty")
			for _, expected := range tt.contains {
				assert.Contains(t, string(content), expected, "migration should contain: %s", expected)
			}
		})
	}
}

func TestMigrationWithSQLite(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_migration_*.db")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

	databaseURL := "sqlite3://" + tmpFile.Name()

	// Run migrations
	err = Run(databaseURL)
	assert.NoError(t, err, "should successfully run migrations")

	// Check version
	version, dirty, err := Version(databaseURL)
	assert.NoError(t, err, "should be able to get version")
	assert.Equal(t, uint(7), version, "should be at version 7")
	assert.False(t, dirty, "should not be in dirty state")
}

func TestMigrationRollback(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_rollback_*.db")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

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

func TestMigrationMultipleRollback(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_multi_rollback_*.db")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

	databaseURL := "sqlite3://" + tmpFile.Name()

	// Run all migrations
	err = Run(databaseURL)
	require.NoError(t, err, "should successfully run migrations")

	// Verify we're at version 7
	version, dirty, err := Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), version, "should be at version 7")
	assert.False(t, dirty)

	// Rollback 3 steps (7 -> 6 -> 5 -> 4)
	err = Rollback(databaseURL, 3)
	assert.NoError(t, err, "should successfully rollback 3 migrations")

	// Verify we're at version 4
	version, dirty, err = Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(4), version, "should be at version 4 after rolling back 3 steps")
	assert.False(t, dirty)

	// Run migrations again to get back to latest
	err = Run(databaseURL)
	assert.NoError(t, err, "should successfully run migrations again")

	// Verify we're back at version 7
	version, dirty, err = Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), version, "should be back at version 7")
	assert.False(t, dirty)
}

func TestMigrationFullCycle(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_full_cycle_*.db")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

	databaseURL := "sqlite3://" + tmpFile.Name()

	// Run all migrations up
	err = Run(databaseURL)
	assert.NoError(t, err, "should successfully run all migrations up")

	version, dirty, err := Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), version, "should be at version 7")
	assert.False(t, dirty)

	// Rollback all migrations (7 steps to get to 0)
	err = Rollback(databaseURL, 7)
	assert.NoError(t, err, "should successfully rollback all migrations")

	// Version should be 0 or return ErrNilVersion
	version, dirty, err = Version(databaseURL)
	if err == nil {
		assert.Equal(t, uint(0), version, "should be at version 0 after full rollback")
		assert.False(t, dirty)
	}

	// Run migrations again
	err = Run(databaseURL)
	assert.NoError(t, err, "should successfully run migrations again after full rollback")

	version, dirty, err = Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), version, "should be back at version 7")
	assert.False(t, dirty)
}

func TestMigrationIdempotency(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "test_idempotency_*.db")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

	databaseURL := "sqlite3://" + tmpFile.Name()

	// Run migrations first time
	err = Run(databaseURL)
	assert.NoError(t, err, "should successfully run migrations first time")

	version, dirty, err := Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), version, "should be at version 7")
	assert.False(t, dirty)

	// Run migrations second time - should be idempotent (no error, no change)
	err = Run(databaseURL)
	assert.NoError(t, err, "running migrations again should not error (idempotent)")

	// Version should still be 7
	version, dirty, err = Version(databaseURL)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), version, "should still be at version 7")
	assert.False(t, dirty)
}

func TestMigrationDownFiles(t *testing.T) {
	// Verify all down migration files are readable and not empty
	downMigrations := []string{
		"000001_initial_schema.down.sql",
		"000002_add_executions_and_indices.down.sql",
		"000003_add_user_created_index.down.sql",
		"000004_add_deleted_at_columns.down.sql",
		"000005_add_max_objs_to_pareto.down.sql",
		"000006_add_updated_at_to_vectors.down.sql",
		"000007_add_execution_metadata.down.sql",
	}

	for _, file := range downMigrations {
		t.Run(file, func(t *testing.T) {
			content, err := migrations.FS.ReadFile(file)
			require.NoError(t, err, "should be able to read down migration file")
			assert.NotEmpty(t, content, "down migration file should not be empty")
			assert.Contains(t, string(content), "DROP", "down migration should contain DROP statements")
		})
	}
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
