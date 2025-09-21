package repository

import (
	"pkg/database"
	"pkg/model"
	"reflect"
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

func TestUserCRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ur := NewSqliteUserRepository(db)

	initialUser := &model.User{
		ID:         "1",
		FirstName:  "fn",
		LastName:   "ln",
		MiddleName: "mn",
	}

	err := ur.CreateOrUpdateUser(initialUser)
	if err != nil {
		t.Fatalf("CreateOrUpdateUser failed on create: %v", err)
	}

	retrievedUser, err := ur.GetUserByID(initialUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if !reflect.DeepEqual(initialUser, retrievedUser) {
		t.Errorf("Retrieved user does not match initial user.\nExpected: %+v\nGot: %+v", initialUser, retrievedUser)
	}

	updatedUser := &model.User{
		ID:         initialUser.ID,
		FirstName:  "newfn",
		LastName:   "newln",
		MiddleName: "newmn",
	}
	err = ur.CreateOrUpdateUser(updatedUser)
	if err != nil {
		t.Fatalf("CreateOrUpdateUser failed on update: %v", err)
	}

	retrievedUpdatedUser, err := ur.GetUserByID(updatedUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed after update: %v", err)
	}
	if !reflect.DeepEqual(updatedUser, retrievedUpdatedUser) {
		t.Errorf("Retrieved user does not match updated user.\nExpected: %+v\nGot: %+v", updatedUser, retrievedUpdatedUser)
	}

	err = ur.DeleteUserByID(updatedUser.ID)
	if err != nil {
		t.Fatalf("DeleteUserByID failed: %v", err)
	}

	deletedUser, err := ur.GetUserByID(updatedUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed when retrieving deleted user: %v", err)
	}
	if deletedUser.ID != "" {
		t.Errorf("Expected user to be deleted, but found user with ID: %s", deletedUser.ID)
	}
}
