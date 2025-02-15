// +build mage

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
	"github.com/magefile/mage/mg"

	"lang_portal_go/internal/models"
	"lang_portal_go/internal/seeder"
)

type DB mg.Namespace

const dbPath = "words.db"

// Init initializes a new SQLite database
func (DB) Init() error {
	fmt.Println("Initializing database...")
	
	// Read schema file
	schema, err := os.ReadFile("db/migrations/001_initial_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// Initialize database using sqlite3
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Execute schema
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	fmt.Println("Database initialized successfully!")
	return nil
}

// Clean removes the database file
func (DB) Clean() error {
	fmt.Println("Cleaning database...")
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove database: %v", err)
	}
	return nil
}

// Seed populates the database with initial data
func (DB) Seed() error {
	fmt.Println("Seeding database...")
	
	db, err := models.NewDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Load and seed initial data
	seedData, err := seeder.LoadSeedData("db/seeds/initial_data.json")
	if err != nil {
		return fmt.Errorf("failed to load seed data: %v", err)
	}

	if err := seeder.SeedDatabase(db, seedData); err != nil {
		return fmt.Errorf("failed to seed database: %v", err)
	}

	// Load and seed test data
	testData, err := seeder.LoadSeedData("db/seeds/test_data.json")
	if err != nil {
		return fmt.Errorf("failed to load test data: %v", err)
	}

	if err := seeder.SeedDatabase(db, testData); err != nil {
		return fmt.Errorf("failed to seed test data: %v", err)
	}

	log.Println("Database seeded successfully!")
	return nil
}

// Reset resets the database to a clean state with seed data
func (DB) Reset() error {
	mg.SerialDeps(DB.Clean, DB.Init, DB.Seed)
	return nil
}

type Server mg.Namespace

// Start starts the application server
func (Server) Start() error {
	fmt.Println("Starting server...")
	cmd := exec.Command("go", "run", "./cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Default target to run when none is specified
var Default = DB.Init
