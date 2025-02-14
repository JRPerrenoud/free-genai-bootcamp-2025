package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "words.db"

func main() {
	fmt.Println("Initializing database...")
	
	// Create database if it doesn't exist
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Read and execute migration file
	migration, err := os.ReadFile("db/migrations/001_initial_schema.sql")
	if err != nil {
		fmt.Printf("Error reading migration file: %v\n", err)
		os.Exit(1)
	}

	_, err = db.Exec(string(migration))
	if err != nil {
		fmt.Printf("Error executing migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully!")
}
