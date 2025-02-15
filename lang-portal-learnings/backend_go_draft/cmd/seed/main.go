package main

import (
	"flag"
	"log"

	"lang_portal_go/internal/models"
	"lang_portal_go/internal/seeder"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Define command line flags
	seedFile := flag.String("file", "db/seeds/initial_data.json", "Path to the seed data JSON file")
	dbFile := flag.String("db", "words.db", "Path to the SQLite database file")
	flag.Parse()

	// Initialize database connection
	db, err := models.NewDB(*dbFile)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Load seed data
	seedData, err := seeder.LoadSeedData(*seedFile)
	if err != nil {
		log.Fatalf("Failed to load seed data: %v", err)
	}

	// Seed the database
	if err := seeder.SeedDatabase(db, seedData); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully!")
}
