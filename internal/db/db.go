package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type InstalledApp struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	InstallPath string   `json:"installPath"`
	Files       []string `json:"files"`
	Shortcut    string   `json:"shortcut"`
	Shortcuts   []ShortcutLink `json:"shortcuts"`
	Startup     bool     `json:"startup"`
	InstalledAt string   `json:"installedAt"`
	ICISource   string   `json:"iciSource"`
	Homepage    string   `json:"homepage"`
}

type ShortcutLink struct {
	Exe  string `json:"exe"`
	Name string `json:"name"`
}

type DB struct {
	conn *sql.DB
	path string
}

func Open() (*DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %w", err)
	}

	dbDir := filepath.Join(configDir, "ICIS")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	dbPath := filepath.Join(dbDir, "icis.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn, path: dbPath}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	var version int
	db.conn.QueryRow("PRAGMA user_version").Scan(&version)

	if version < 1 {
		_, err := db.conn.Exec(`
			CREATE TABLE IF NOT EXISTS installs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL UNIQUE,
				version TEXT NOT NULL DEFAULT '',
				install_path TEXT NOT NULL,
				files TEXT NOT NULL DEFAULT '[]',
				shortcut TEXT NOT NULL DEFAULT '',
				startup INTEGER NOT NULL DEFAULT 0,
				installed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				ici_source TEXT NOT NULL DEFAULT ''
			)
		`)
		if err != nil {
			return err
		}
		db.conn.Exec("PRAGMA user_version = 1")
	}

	if version < 2 {
		_, err := db.conn.Exec(`ALTER TABLE installs ADD COLUMN shortcuts TEXT NOT NULL DEFAULT '[]'`)
		if err != nil {
			return err
		}
		db.conn.Exec("PRAGMA user_version = 2")
	}

	if version < 3 {
		_, err := db.conn.Exec(`ALTER TABLE installs ADD COLUMN homepage TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			return err
		}
		db.conn.Exec("PRAGMA user_version = 3")
	}

	return nil
}

func (db *DB) SaveApp(app InstalledApp) error {
	filesJSON, err := json.Marshal(app.Files)
	if err != nil {
		return err
	}

	shortcutsJSON, err := json.Marshal(app.Shortcuts)
	if err != nil {
		return err
	}

	startupInt := 0
	if app.Startup {
		startupInt = 1
	}

	if app.InstalledAt == "" {
		app.InstalledAt = time.Now().Format(time.RFC3339)
	}

	_, err = db.conn.Exec(`
		INSERT OR REPLACE INTO installs (name, version, install_path, files, shortcut, shortcuts, startup, installed_at, ici_source, homepage)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, app.Name, app.Version, app.InstallPath, string(filesJSON), app.Shortcut, string(shortcutsJSON), startupInt, app.InstalledAt, app.ICISource, app.Homepage)

	return err
}

func (db *DB) GetApp(name string) (*InstalledApp, error) {
	row := db.conn.QueryRow(`
		SELECT id, name, version, install_path, files, shortcut, shortcuts, startup, installed_at, ici_source, homepage
		FROM installs WHERE name = ?
	`, name)

	var app InstalledApp
	var filesJSON string
	var shortcutsJSON string
	var startupInt int

	err := row.Scan(&app.ID, &app.Name, &app.Version, &app.InstallPath, &filesJSON, &app.Shortcut, &shortcutsJSON, &startupInt, &app.InstalledAt, &app.ICISource, &app.Homepage)
	if err != nil {
		return nil, err
	}

	app.Startup = startupInt == 1
	if err := json.Unmarshal([]byte(filesJSON), &app.Files); err != nil {
		app.Files = []string{}
	}
	if err := json.Unmarshal([]byte(shortcutsJSON), &app.Shortcuts); err != nil {
		app.Shortcuts = []ShortcutLink{}
	}

	return &app, nil
}

func (db *DB) ListApps() ([]InstalledApp, error) {
	rows, err := db.conn.Query(`
		SELECT id, name, version, install_path, files, shortcut, shortcuts, startup, installed_at, ici_source, homepage
		FROM installs ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []InstalledApp
	for rows.Next() {
		var app InstalledApp
		var filesJSON string
		var shortcutsJSON string
		var startupInt int

		if err := rows.Scan(&app.ID, &app.Name, &app.Version, &app.InstallPath, &filesJSON, &app.Shortcut, &shortcutsJSON, &startupInt, &app.InstalledAt, &app.ICISource, &app.Homepage); err != nil {
			return nil, err
		}

		app.Startup = startupInt == 1
		if err := json.Unmarshal([]byte(filesJSON), &app.Files); err != nil {
			app.Files = []string{}
		}
		if err := json.Unmarshal([]byte(shortcutsJSON), &app.Shortcuts); err != nil {
			app.Shortcuts = []ShortcutLink{}
		}

		apps = append(apps, app)
	}

	return apps, rows.Err()
}

func (db *DB) DeleteApp(name string) error {
	_, err := db.conn.Exec("DELETE FROM installs WHERE name = ?", name)
	return err
}
