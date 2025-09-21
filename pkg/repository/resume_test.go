package repository

import (
	"pkg/database"
	"pkg/model"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*database.DB, func()) {
	db, err := database.NewSqliteDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize in-memory database: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}
func TestUserCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ur := NewSqliteUserRepository(db)

	id := "1"
	fn := "fn"
	ln := "ln"
	mn := "mn"

	err := ur.CreateOrUpdateUser(&model.User{ID: id, FirstName: fn, LastName: ln, MiddleName: mn})
	if err != nil {
		t.Errorf("failed creating or updating user %v", err)
	}

	var u model.User
	err = db.QueryRow("select id, first_name, last_name, middle_name from users where id = ?", "1").Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.MiddleName,
	)
	if err != nil {
		t.Errorf("failed quering row %v", err)
	}

	if u.ID != id {
		t.Errorf("Expected ID to be equal '%s', got '%s'", id, u.ID)
	}
	if u.FirstName != fn {
		t.Errorf("Expected FirstName to be equal '%s', got '%s'", fn, u.FirstName)
	}
	if u.LastName != ln {
		t.Errorf("Expected LastName to be equal '%s', got '%s'", ln, u.LastName)
	}
	if u.MiddleName != mn {
		t.Errorf("Expected MiddleName to be equal '%s', got '%s'", mn, u.MiddleName)
	}
}

func TestUserUpdate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ur := NewSqliteUserRepository(db)

	id := "1"
	fn := "fn"
	ln := "ln"
	mn := "mn"

	err := ur.CreateOrUpdateUser(&model.User{ID: id, FirstName: fn, LastName: ln, MiddleName: mn})
	if err != nil {
		t.Errorf("failed creating or updating user %v", err)
	}

	newfn := "newfn"
	newln := "newln"
	newmn := "newmn"
	err = ur.CreateOrUpdateUser(&model.User{ID: id, FirstName: newfn, LastName: newln, MiddleName: newmn})
	if err != nil {
		t.Errorf("failed creating or updating user %v", err)
	}

	var u model.User
	err = db.QueryRow("select id, first_name, last_name, middle_name from users where id = ?", "1").Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.MiddleName,
	)
	if err != nil {
		t.Errorf("failed quering row %v", err)
	}

	if u.ID != id {
		t.Errorf("Expected ID to be equal '%s', got '%s'", id, u.ID)
	}
	if u.FirstName != newfn {
		t.Errorf("Expected FirstName to be equal '%s', got '%s'", newfn, u.FirstName)
	}
	if u.LastName != newln {
		t.Errorf("Expected LastName to be equal '%s', got '%s'", newln, u.LastName)
	}
	if u.MiddleName != newmn {
		t.Errorf("Expected MiddleName to be equal '%s', got '%s'", newmn, u.MiddleName)
	}
}
