package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func main() {
	// Get the project root directory by going up two levels from the current file
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	dbPath := filepath.Join(projectRoot, "words.db")

	// Create database directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Printf("Error creating database directory: %v\n", err)
		os.Exit(1)
	}

	// Read migration and seed files
	migrationSQL, err := os.ReadFile(filepath.Join(projectRoot, "db", "migrations", "000001_create_initial_schema.up.sql"))
	if err != nil {
		fmt.Printf("Error reading migration file: %v\n", err)
		os.Exit(1)
	}

	seedSQL, err := os.ReadFile(filepath.Join(projectRoot, "db", "seeds", "000001_initial_seed.sql"))
	if err != nil {
		fmt.Printf("Error reading seed file: %v\n", err)
		os.Exit(1)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Execute migration SQL
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		fmt.Printf("Error executing migrations: %v\n", err)
		os.Exit(1)
	}

	// Execute seed SQL
	_, err = db.Exec(string(seedSQL))
	if err != nil {
		fmt.Printf("Error executing seeds: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully!")
}
