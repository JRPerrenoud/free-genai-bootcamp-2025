// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Build builds the application
func Build() error {
	return sh.Run("go", "build", "./cmd/server")
}

// Run runs the server
func Run() error {
	return sh.Run("go", "run", "./cmd/server")
}

// InitDB initializes the database
func InitDB() error {
	return sh.Run("go", "run", "./cmd/initdb")
}

// Seed seeds the database with initial data
func Seed() error {
	return sh.Run("go", "run", "./cmd/seed")
}

// Test runs the test suite
func Test() error {
	return sh.Run("go", "test", "./...")
}

// Clean cleans build artifacts
func Clean() error {
	return sh.Run("rm", "-f", "backend_go")
}
