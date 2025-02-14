package main

import (
	"flag"
	"log"
	"os"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Define command line flags
	dbFile := flag.String("db", "words.db", "Path to the SQLite database file")
	schemaFile := flag.String("schema", "db/migrations/001_initial_schema.sql", "Path to the schema file")
	flag.Parse()

	// Read schema file
	schema, err := os.ReadFile(*schemaFile)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	// Remove existing database if it exists
	os.Remove(*dbFile)

	// Create and open new database
	db, err := sql.Open("sqlite3", *dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Execute schema
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
	}

	log.Println("Database initialized successfully!")
}
