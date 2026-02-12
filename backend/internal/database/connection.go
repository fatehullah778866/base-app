package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"base-app-service/migrations"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
	logger *zap.Logger
}

type DatabaseConfig struct {
	Driver                string
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	SSLMode               string
	SQLitePath            string
	MaxConnections        int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

func NewConnection(cfg DatabaseConfig, logger *zap.Logger) (*DB, error) {
	var (
		driver = cfg.Driver
		dsn    string
	)

	if driver == "" {
		driver = "sqlite"
	}

	switch driver {
	case "sqlite":
		dsn = cfg.SQLitePath
		if dsn == "" {
			if cfg.Name != "" {
				dsn = fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)", cfg.Name)
			} else {
				dsn = "file:app.db?_pragma=foreign_keys(ON)"
			}
		}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool settings
	if driver == "sqlite" {
		// SQLite needs a small pool; shared cache is handled by the driver DSN
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	} else {
		db.SetMaxOpenConns(cfg.MaxConnections)
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
		db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")

	return &DB{DB: db, logger: logger}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Transaction helper
func (db *DB) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// RunMigrations executes all .sql files in the given directory in lexical order.
// This keeps migration tooling simple for the SQLite setup.
func (db *DB) RunMigrations(migrationsDir string) error {
	// Track applied migrations to avoid re-running ALTER statements
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations_applied (name TEXT PRIMARY KEY, applied_at DATETIME DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		return fmt.Errorf("create migrations_applied: %w", err)
	}

	// Ensure critical columns/tables exist even if older DB file is present
	if err := db.ensureAdminSupport(); err != nil {
		return err
	}

	migrationFiles, err := listMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("collect migrations: %w", err)
	}

	for _, mf := range migrationFiles {
		filename := mf.name
		var count int
		if err := db.QueryRow(`SELECT COUNT(1) FROM migrations_applied WHERE name = ?`, filename).Scan(&count); err != nil {
			return fmt.Errorf("check migration %s: %w", filename, err)
		}
		if count > 0 {
			continue // already applied
		}

		if _, err := db.Exec(string(mf.content)); err != nil {
			// Skip benign "already exists" errors for idempotency
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "duplicate column name") || strings.Contains(msg, "already exists") {
				_, _ = db.Exec(`INSERT OR IGNORE INTO migrations_applied (name) VALUES (?)`, filename)
				continue
			}
			return fmt.Errorf("exec migration %s: %w", filename, err)
		}

		if _, err := db.Exec(`INSERT INTO migrations_applied (name) VALUES (?)`, filename); err != nil {
			return fmt.Errorf("record migration %s: %w", filename, err)
		}
	}

	return nil
}

// ensureAdminSupport patches older SQLite files by adding the role column and admin tables if missing.
func (db *DB) ensureAdminSupport() error {
	// Add role column if missing
	type colInfo struct {
		name string
	}
	var hasRole bool
	rows, err := db.Query(`PRAGMA table_info('users')`)
	if err != nil {
		return fmt.Errorf("inspect users schema: %w", err)
	}
	defer rows.Close()
	var foundAny bool
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err == nil {
			foundAny = true
			if strings.EqualFold(name, "role") {
				hasRole = true
			}
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	// If the users table doesn't exist yet, skip attempting to ALTER it
	if !foundAny {
		return nil
	}

	if !hasRole {
		if _, err := db.Exec(`ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'user'`); err != nil {
			// ignore if concurrently added
			if !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
				return fmt.Errorf("add users.role: %w", err)
			}
		}
		_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`)
	}

	// Create admin-related tables if missing
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS activity_logs (
		id TEXT PRIMARY KEY,
		actor_id TEXT,
		actor_role TEXT,
		action TEXT NOT NULL,
		target_type TEXT,
		target_id TEXT,
		metadata TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_activity_logs_actor ON activity_logs(actor_id);
	CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
	CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at);
	`); err != nil {
		return fmt.Errorf("ensure activity_logs: %w", err)
	}

	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS access_requests (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT,
		details TEXT,
		status TEXT DEFAULT 'pending',
		feedback TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_access_requests_status ON access_requests(status);
	CREATE INDEX IF NOT EXISTS idx_access_requests_user ON access_requests(user_id);
	`); err != nil {
		return fmt.Errorf("ensure access_requests: %w", err)
	}

	return nil
}

type migrationFile struct {
	name    string
	content []byte
}

func listMigrationFiles(dir string) ([]migrationFile, error) {
	fileMap := map[string][]byte{}

	if dir != "" {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
					continue
				}
				content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
				if err != nil {
					return nil, fmt.Errorf("read migration %s: %w", entry.Name(), err)
				}
				fileMap[entry.Name()] = content
			}
		} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("read migrations dir: %w", err)
		}
	}

	if entries, err := migrations.UpFS.ReadDir("."); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
				continue
			}
			if _, exists := fileMap[entry.Name()]; exists {
				continue
			}
			content, err := migrations.UpFS.ReadFile(entry.Name())
			if err != nil {
				return nil, fmt.Errorf("read embedded migration %s: %w", entry.Name(), err)
			}
			fileMap[entry.Name()] = content
		}
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}

	var migrationsSlice []migrationFile
	for name, content := range fileMap {
		migrationsSlice = append(migrationsSlice, migrationFile{name: name, content: content})
	}
	sort.Slice(migrationsSlice, func(i, j int) bool {
		return migrationsSlice[i].name < migrationsSlice[j].name
	})

	return migrationsSlice, nil
}
