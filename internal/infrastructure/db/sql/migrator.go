package dbsql

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultMigrationsDir = "migrations"

type migrationFile struct {
	version string
	path    string
}

// RunMigrations は指定されたディレクトリ内の .sql ファイルを実行し、schema_migrations に記録する。
func RunMigrations(ctx context.Context, db *sql.DB, dir string) error {
	if dir == "" {
		dir = defaultMigrationsDir
	}

	files, err := loadMigrationFiles(dir)
	if err != nil {
		return err
	}

	if err := ensureSchemaMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	for _, mf := range files {
		applied, err := isApplied(ctx, db, mf.version)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", mf.version, err)
		}
		if applied {
			continue
		}

		if err := applyMigration(ctx, db, mf); err != nil {
			return fmt.Errorf("apply migration %s: %w", mf.version, err)
		}
		log.Printf("[migrate] applied %s", mf.version)
	}

	return nil
}

func loadMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	var files []migrationFile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, migrationFile{
			version: entry.Name(),
			path:    filepath.Join(dir, entry.Name()),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].version < files[j].version
	})

	return files, nil
}

func ensureSchemaMigrationsTable(ctx context.Context, db *sql.DB) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL
)`
	_, err := db.ExecContext(ctx, query)
	return err
}

func isApplied(ctx context.Context, db *sql.DB, version string) (bool, error) {
	const query = `SELECT 1 FROM schema_migrations WHERE version = $1`
	var exists int
	err := db.QueryRowContext(ctx, query, version).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func applyMigration(ctx context.Context, db *sql.DB, mf migrationFile) error {
	f, err := os.Open(mf.path)
	if err != nil {
		return fmt.Errorf("open migration file %s: %w", mf.path, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read migration file %s: %w", mf.path, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf("exec migration %s: %w", mf.version, err)
	}

	const insert = `INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)`
	if _, err = tx.ExecContext(ctx, insert, mf.version, time.Now().UTC()); err != nil {
		return fmt.Errorf("mark migration %s: %w", mf.version, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", mf.version, err)
	}

	return nil
}
