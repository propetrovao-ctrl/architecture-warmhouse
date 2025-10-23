package migrations

import (
	"context"
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var migrationFiles embed.FS

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	Filename  string
	Content   string
	Checksum  string
}

// Migrator handles database migrations
type Migrator struct {
	pool *pgxpool.Pool
}

// NewMigrator creates a new migrator instance
func NewMigrator(pool *pgxpool.Pool) *Migrator {
	return &Migrator{pool: pool}
}

// RunMigrations executes all pending migrations
func (m *Migrator) RunMigrations(ctx context.Context) error {
	// Ensure migrations table exists
	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find pending migrations
	pending := m.findPendingMigrations(migrations, applied)

	if len(pending) == 0 {
		fmt.Println("No pending migrations found")
		return nil
	}

	// Execute pending migrations
	for _, migration := range pending {
		fmt.Printf("Running migration: %s\n", migration.Version)
		if err := m.executeMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}
		fmt.Printf("Migration %s completed successfully\n", migration.Version)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(50) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			checksum VARCHAR(64)
		)
	`
	_, err := m.pool.Exec(ctx, query)
	return err
}

// loadMigrations loads all migration files from the embedded filesystem
func (m *Migrator) loadMigrations() ([]Migration, error) {
	var migrations []Migration

	err := fs.WalkDir(migrationFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		content, err := migrationFiles.ReadFile(path)
		if err != nil {
			return err
		}

		// Extract version and name from filename
		filename := filepath.Base(path)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename format: %s", filename)
		}

		version := parts[0]
		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")

		// Calculate checksum
		checksum := fmt.Sprintf("%x", sha256.Sum256(content))

		migration := Migration{
			Version:  version,
			Name:     name,
			Filename: filename,
			Content:  string(content),
			Checksum: checksum,
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a map of applied migration versions
func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	query := "SELECT version FROM migrations ORDER BY version"
	rows, err := m.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// findPendingMigrations finds migrations that haven't been applied yet
func (m *Migrator) findPendingMigrations(migrations []Migration, applied map[string]bool) []Migration {
	var pending []Migration
	for _, migration := range migrations {
		if !applied[migration.Version] {
			pending = append(pending, migration)
		}
	}
	return pending
}

// executeMigration executes a single migration
func (m *Migrator) executeMigration(ctx context.Context, migration Migration) error {
	// Start transaction
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Execute migration SQL
	_, err = tx.Exec(ctx, migration.Content)
	if err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	// Record migration as applied
	recordQuery := `
		INSERT INTO migrations (version, name, applied_at, checksum)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, recordQuery, migration.Version, migration.Name, time.Now(), migration.Checksum)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	return tx.Commit(ctx)
}

// GetMigrationStatus returns the status of all migrations
func (m *Migrator) GetMigrationStatus(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT version, name, applied_at, checksum
		FROM migrations
		ORDER BY version
	`
	rows, err := m.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var status []map[string]interface{}
	for rows.Next() {
		var version, name, checksum string
		var appliedAt time.Time

		if err := rows.Scan(&version, &name, &appliedAt, &checksum); err != nil {
			return nil, err
		}

		status = append(status, map[string]interface{}{
			"version":    version,
			"name":       name,
			"applied_at": appliedAt,
			"checksum":   checksum,
		})
	}

	return status, rows.Err()
}
