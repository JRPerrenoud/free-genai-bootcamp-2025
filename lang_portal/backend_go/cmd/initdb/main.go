package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "db/lang_portal.db"
	}

	// Get absolute path
	absDbPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to database: %v", err)
	}

	// Remove existing database if it exists
	if _, err := os.Stat(absDbPath); err == nil {
		if err := os.Remove(absDbPath); err != nil {
			log.Fatalf("Failed to remove existing database: %v", err)
		}
	}

	// Create new database
	db, err := sql.Open("sqlite3", absDbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Read and execute all migration files in order
	migrationsDir := filepath.Join("db", "migrations")
	migrations, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Execute each migration file
	for _, migration := range migrations {
		if !migration.IsDir() && filepath.Ext(migration.Name()) == ".sql" {
			migrationPath := filepath.Join(migrationsDir, migration.Name())
			migrationSQL, err := os.ReadFile(migrationPath)
			if err != nil {
				log.Fatalf("Failed to read migration file %s: %v", migration.Name(), err)
			}

			// Execute migration
			_, err = db.Exec(string(migrationSQL))
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v", migration.Name(), err)
			}
			fmt.Printf("Applied migration: %s\n", migration.Name())
		}
	}

	fmt.Println("Database initialized successfully!")
}
