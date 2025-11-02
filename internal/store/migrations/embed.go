// Package migrations contains embedded database migration files.
package migrations

import "embed"

// FS contains all migration SQL files embedded into the binary.
//
//go:embed *.sql
var FS embed.FS
