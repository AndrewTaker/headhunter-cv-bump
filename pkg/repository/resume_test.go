package repository

import (
	"database/sql"
	"log"
	"pkg/database"
	"pkg/model"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var userID = "user-123"
var fn = "fn"
var ln = "ln"
var mn = "mn"
var ca1 = "2023-05-25T14:19:02+0300"
var ua1 = "2023-05-25T14:20:02+0300"
var ca2 = "2023-05-25T14:19:02+0300"
var ua2 = "2023-05-25T14:20:02+0300"

func hht(hhtime string) model.HHTime {
	t, _ := time.Parse(model.TimeLayout, hhtime)
	return model.HHTime(t)
}

func TestQQ(t *testing.T) {
	db, cleanup := setupResumeTestDB(t)
	defer cleanup()
	rr := NewSqliteResumeRepository(db)
	r := []model.Resume{
		{
			ID:           "r1",
			Title:        "r1-title",
			AlternateUrl: "r1-url",
			CreatedAt:    hht(ca1),
			UpdatedAt:    hht(ua1),
		},
	}
	err := rr.CreateOrUpdateResumes(r, userID)
	if err != nil {
		t.Fatalf("TestQQ failed %v", err)
	}
	r0 := r[0]
	r0t := time.Time(r0.CreatedAt)
	log.Println(r0t.Format(time.RFC3339))

}

func setupResumeTestDB(t *testing.T) (*database.DB, func()) {
	db, err := database.NewSqliteDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize in-memory database: %v", err)
	}

	ur := NewSqliteUserRepository(db)
	if err := ur.CreateOrUpdateUser(&model.User{ID: userID, FirstName: fn, LastName: ln, MiddleName: mn}); err != nil {
		t.Fatalf("Failed to create test user %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestHHTime_Scan(t *testing.T) {
	var hht model.HHTime

	testTime := time.Now()
	err := hht.Scan(testTime)
	if err != nil {
		t.Errorf("Scan failed for time.Time: %v", err)
	}
	if !time.Time(hht).Equal(testTime) {
		t.Errorf("Scan from time.Time failed. Expected %v, got %v", testTime, time.Time(hht))
	}

	testTimeStr := "2023-10-27 10:30:00+00:00"
	tt, _ := time.Parse(model.TimeLayout, testTimeStr)
	err = hht.Scan(testTimeStr)
	if err != nil {
		t.Errorf("Scan failed for string: %v", err)
	}
	if !time.Time(hht).Equal(tt) {
		t.Errorf("Scan from string failed. Expected %v, got %v", t, time.Time(hht))
	}

	testTimeBytes := []byte("2023-10-27 11:30:00+00:00")
	tt, _ = time.Parse(model.TimeLayout, string(testTimeBytes))
	err = hht.Scan(testTimeBytes)
	if err != nil {
		t.Errorf("Scan failed for []byte: %v", err)
	}
	if !time.Time(hht).Equal(tt) {
		t.Errorf("Scan from []byte failed. Expected %v, got %v", t, time.Time(hht))
	}
}

func TestResumeRepositoryCRUD(t *testing.T) {
	db, cleanup := setupResumeTestDB(t)
	defer cleanup()
	rr := NewSqliteResumeRepository(db)

	t.Run("Create and Get Resumes", func(t *testing.T) {
		resumes := []model.Resume{
			{
				ID:           "resume-1",
				Title:        "Resume A",
				AlternateURL: "http://example.com/a",
				CreatedAt:    hht(ca1),
				UpdatedAt:    hht(ua1),
			},
			{
				ID:           "resume-2",
				Title:        "Resume B",
				AlternateURL: "http://example.com/b",
				CreatedAt:    hht(ca2),
				UpdatedAt:    hht(ua2),
			},
		}

		err := rr.CreateOrUpdateResumes(resumes, userID)
		if err != nil {
			t.Fatalf("CreateOrUpdateResumes failed: %v", err)
		}

		retrievedResumes, err := rr.GetUserResumes(userID)
		if err != nil {
			t.Fatalf("GetUserResumes failed: %v", err)
		}
		if len(retrievedResumes) != 2 {
			t.Errorf("Expected 2 resumes, got %d", len(retrievedResumes))
		}

		retrievedMap := make(map[string]model.Resume)
		for _, r := range retrievedResumes {
			retrievedMap[r.ID] = r
		}
		for _, r := range resumes {
			if r.CreatedAt != retrievedMap[r.ID].CreatedAt {
				t.Errorf("%v+ != %v+", r.CreatedAt, retrievedMap[r.ID].CreatedAt)
			}
		}
	})

	t.Run("Update an existing resume", func(t *testing.T) {
		updatedResume := []model.Resume{
			{
				ID:           "resume-1",
				Title:        "Updated Resume A",
				AlternateURL: "http://example.com/a-updated",
				CreatedAt:    hht(ca2),
				UpdatedAt:    hht(ua2),
			},
		}
		err := rr.CreateOrUpdateResumes(updatedResume, userID)
		if err != nil {
			t.Fatalf("CreateOrUpdateResumes failed on update: %v", err)
		}

		retrievedUpdatedResumes, err := rr.GetUserResumes(userID)
		if err != nil {
			t.Fatalf("GetUserResumes failed after update: %v", err)
		}
		var updatedItem *model.Resume
		for i := range retrievedUpdatedResumes {
			if retrievedUpdatedResumes[i].ID == "resume-1" {
				updatedItem = &retrievedUpdatedResumes[i]
				break
			}
		}
		if updatedItem == nil {
			t.Fatal("Updated resume not found")
		}
		if updatedItem.Title != "Updated Resume A" {
			t.Errorf("Expected updated title 'Updated Resume A', got '%s'", updatedItem.Title)
		}
		if updatedItem.AlternateURL != "http://example.com/a-updated" {
			t.Errorf("Expected updated URL, got '%s'", updatedItem.AlternateURL)
		}
	})

	t.Run("Get a single resume by ID", func(t *testing.T) {
		otherUserID := "user-456"
		rr.CreateOrUpdateResumes([]model.Resume{{ID: "resume-3", Title: "Resume 3"}}, otherUserID)

		retrievedResume, err := rr.GetResumeByID("resume-1", userID)
		if err != nil {
			t.Fatalf("GetResumeByID failed for an existing resume: %v", err)
		}
		if retrievedResume.ID != "resume-1" {
			t.Errorf("Expected resume-1, got %s", retrievedResume.ID)
		}

		nonExistentResume, err := rr.GetResumeByID("non-existent-id", userID)
		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, but got %v", err)
		}
		if nonExistentResume.ID != "" {
			t.Errorf("Expected an empty user, got %s", nonExistentResume.ID)
		}

		wrongUserResume, err := rr.GetResumeByID("resume-1", otherUserID)
		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows for wrong user, but got %v", err)
		}
		if wrongUserResume.ID != "" {
			t.Errorf("Expected an empty user, got %s", wrongUserResume.ID)
		}
	})

	t.Run("Toggle Scheduling", func(t *testing.T) {
		resumeID := "resume-to-toggle"
		rr.CreateOrUpdateResumes([]model.Resume{{ID: resumeID, Title: "Toggle Me"}}, userID)

		err := rr.ToggleScheduling(resumeID, userID, true)
		if err != nil {
			t.Fatalf("ToggleScheduling failed to set to true: %v", err)
		}
		retrieved, err := rr.GetResumeByID(resumeID, userID)
		if err != nil {
			t.Fatalf("GetResumeByID failed: %v", err)
		}
		if retrieved.IsScheduled != 1 {
			t.Errorf("Expected is_scheduled to be 1, got %d", retrieved.IsScheduled)
		}

		err = rr.ToggleScheduling(resumeID, userID, false)
		if err != nil {
			t.Fatalf("ToggleScheduling failed to set to false: %v", err)
		}
		retrieved, err = rr.GetResumeByID(resumeID, userID)
		if err != nil {
			t.Fatalf("GetResumeByID failed: %v", err)
		}
		if retrieved.IsScheduled != 0 {
			t.Errorf("Expected is_scheduled to be 0, got %d", retrieved.IsScheduled)
		}
	})

	t.Run("Delete resumes", func(t *testing.T) {
		otherUserID := "user-456"
		resumesToDelete := []model.Resume{{ID: "resume-1", Title: "A"}, {ID: "resume-2", Title: "B"}}
		resumesToKeep := []model.Resume{{ID: "resume-3", Title: "C"}}
		rr.CreateOrUpdateResumes(resumesToDelete, userID)
		rr.CreateOrUpdateResumes(resumesToKeep, userID)
		rr.CreateOrUpdateResumes([]model.Resume{{ID: "resume-4", Title: "D"}}, otherUserID)

		err := rr.DeleteResumesByUserID(resumesToDelete, userID)
		if err != nil {
			t.Fatalf("DeleteResumesByUserID failed: %v", err)
		}

		retrieved, err := rr.GetUserResumes(userID)
		if err != nil {
			t.Fatalf("GetUserResumes failed after deletion: %v", err)
		}
		if len(retrieved) != 1 || retrieved[0].ID != "resume-3" {
			t.Fatalf("Expected 1 resume after deletion (ID 'resume-3'), got %d", len(retrieved))
		}
		otherUserResumes, err := rr.GetUserResumes(otherUserID)
		if err != nil {
			t.Fatalf("GetUserResumes failed for other user: %v", err)
		}
		if len(otherUserResumes) != 1 || otherUserResumes[0].ID != "resume-4" {
			t.Errorf("Expected other user's resume to be untouched, but was affected")
		}
	})
}
