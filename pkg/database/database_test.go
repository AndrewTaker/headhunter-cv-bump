package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	db, err := NewSqliteDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize in-memory database: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestTableCreatedUsers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	var tableName string
	err := db.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		t.Errorf("Expected 'users' table to exist, but got error: %v", err)
	}
	if tableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", tableName)
	}
}

func TestTableCreatedResumes(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	var tableName string
	err := db.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='resumes'").Scan(&tableName)
	if err != nil {
		t.Errorf("Expected 'resumes' table to exist, but got error: %v", err)
	}
	if tableName != "resumes" {
		t.Errorf("Expected table name 'resumes', got '%s'", tableName)
	}
}

func TestTableCreatedTokens(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	var tableName string
	err := db.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tokens'").Scan(&tableName)
	if err != nil {
		t.Errorf("Expected 'tokens' table to exist, but got error: %v", err)
	}
	if tableName != "tokens" {
		t.Errorf("Expected table name 'tokens', got '%s'", tableName)
	}
}

func TestTableCreatedScheduler(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	var tableName string
	err := db.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler'").Scan(&tableName)
	if err != nil {
		t.Errorf("Expected 'scheduler' table to exist, but got error: %v", err)
	}
	if tableName != "scheduler" {
		t.Errorf("Expected table name 'scheduler', got '%s'", tableName)
	}
}
