package infrastructure

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes all pending migrations
func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
	// Create schema_migrations table
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Read migration files
	migrations, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sort.Strings(migrations)

	// Apply each migration
	for _, migration := range migrations {
		version, err := extractVersion(migration)
		if err != nil {
			return fmt.Errorf("failed to extract version from %s: %w", migration, err)
		}

		// Check if already applied
		var exists bool
		err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)", version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			continue
		}

		// Read and execute migration
		sql, err := os.ReadFile(migration)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		_, err = db.Exec(ctx, string(sql))
		if err != nil {
			return fmt.Errorf("migration %s failed: %w", migration, err)
		}

		// Record migration
		_, err = db.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
		if err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		fmt.Printf("Applied migration: %s\n", filepath.Base(migration))
	}

	return nil
}

func extractVersion(filename string) (int, error) {
	base := filepath.Base(filename)
	parts := strings.SplitN(base, "_", 2)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid version number: %s", parts[0])
	}

	return version, nil
}
