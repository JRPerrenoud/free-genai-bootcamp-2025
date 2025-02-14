//go:build mage
// +build mage

// Package main provides mage build targets for the project
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/magefile/mage/sh"
)

const (
	dbPath = "./db/lang_portal.db"
)

// Default target to run when none is specified
var Default = Build

type DB mg.Namespace

// Build builds the application
func Build() error {
	fmt.Println("Building...")
	return sh.RunV("go", "build", "-o", "bin/server", "./cmd/server")
}

// Run runs the application
func Run() error {
	mg.Deps(Build)
	fmt.Println("Running server...")
	return sh.RunV("./bin/server")
}

// Clean cleans build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	return os.RemoveAll("bin")
}

// Test runs the tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "./...")
}

// Deps installs dependencies
func Deps() error {
	fmt.Println("Installing dependencies...")
	return sh.RunV("go", "mod", "download")
}

// Init initializes the database
func (DB) Init() error {
	fmt.Println("Initializing database...")

	// Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("error creating database directory: %v", err)
	}

	// Remove existing database if it exists
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing existing database: %v", err)
	}

	// Create new database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Execute schema migrations
	schemaSQL, err := os.ReadFile("./db/migrations/000001_create_initial_schema.up.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	if _, err := db.Exec(string(schemaSQL)); err != nil {
		return fmt.Errorf("error executing migrations: %v", err)
	}

	// Execute seeds
	seedSQL, err := os.ReadFile("./db/seeds/000001_initial_data.sql")
	if err != nil {
		return fmt.Errorf("error reading seed file: %v", err)
	}

	if _, err := db.Exec(string(seedSQL)); err != nil {
		return fmt.Errorf("error executing seeds: %v", err)
	}

	fmt.Println("Database initialized successfully!")
	return nil
}
