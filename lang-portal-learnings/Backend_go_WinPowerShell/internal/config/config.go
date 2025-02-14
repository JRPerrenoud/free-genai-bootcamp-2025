package config

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Config holds all configuration for our application
type Config struct {
	DB     *sql.DB
	Port   string
	DBPath string
}

// New creates a new Config
func New() (*Config, error) {
	cfg := &Config{
		Port:   "8081",
		DBPath: "words.db",
	}

	// Initialize database connection
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	cfg.DB = db
	return cfg, nil
}

// Close closes the database connection
func (c *Config) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
