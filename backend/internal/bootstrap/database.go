package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Database initialise la connexion à la base de données et exécute les migrations
func Database(dataDir string) (*sql.DB, error) {
	if dataDir == "" {
		dataDir = "/data"
	}

	// Créer les répertoires
	mediaDir := filepath.Join(dataDir, "media")
	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create media dir: %w", err)
	}

	// Ouvrir la base de données
	dbPath := filepath.Join(dataDir, "db.sqlite")
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_fk=1", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Exécuter les migrations
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	log.Println("database initialized successfully")
	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		password TEXT,
		is_active INTEGER DEFAULT 0,
		is_admin INTEGER DEFAULT 0,
		created_at INTEGER
	);
	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		owner_id TEXT,
		filename TEXT,
		size INTEGER,
		path TEXT,
		created_at INTEGER
	);
	`
	_, err := db.Exec(schema)
	return err
}

// GetMediaDir retourne le chemin du répertoire média
func GetMediaDir(dataDir string) string {
	if dataDir == "" {
		dataDir = "/data"
	}
	return filepath.Join(dataDir, "media")
}
