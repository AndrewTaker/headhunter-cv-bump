package repository

import (
	"database/sql"
	"fmt"
	"pkg/database"
	"pkg/model"
)

type ResumeRepository interface {
	CreateOrUpdateResumes(resumes []model.Resume, userID string) error
	GetResumesByUserID(userID string) ([]model.Resume, error)
	GetResumeByID(resumeID, userID string) (*model.Resume, error)
	ToggleScheduling(resumeID, userID string, isScheduled bool) error
}

type SqliteResumeRepository struct {
	DB *database.DB
}

func NewSqliteResumeRepository(db *database.DB) ResumeRepository {
	return &SqliteResumeRepository{DB: db}
}

func (ur *SqliteResumeRepository) CreateOrUpdateResumes(resumes []model.Resume, userID string) error {
	query := `
	insert into resumes (id, title, alternate_url, created_at, updated_at, user_id) values (?, ?, ?, ?, ?, ?)
	on conflict(id) do update set
	title = excluded.title,
	alternate_url = excluded.alternate_url,
	created_at = excluded.created_at,
	updated_at = excluded.updated_at,
	user_id = excluded.user_id;
	`
	tx, err := ur.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, resume := range resumes {
		_, err := stmt.Exec(
			resume.ID,
			resume.Title,
			resume.AlternateURL,
			resume.CreatedAt,
			resume.UpdatedAt,
			userID,
		)
		if err != nil {
			return fmt.Errorf("failed to execute statement for resume ID %s: %w", resume.ID, err)
		}
	}
	tx.Commit()

	return nil
}

func (ur *SqliteResumeRepository) GetResumesByUserID(userID string) ([]model.Resume, error) {
	query := "select id, title, alternate_url, created_at, updated_at, is_scheduled from resumes where user_id = ?"

	rows, err := ur.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	var resumes []model.Resume
	for rows.Next() {
		var r model.Resume
		if err := rows.Scan(
			&r.ID,
			&r.Title,
			&r.AlternateURL,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.IsScheduled,
		); err != nil {
			return nil, err
		}
		resumes = append(resumes, r)
	}

	if err = rows.Err(); err != nil {
		return []model.Resume{}, err
	}

	if err == sql.ErrNoRows {
		return []model.Resume{}, nil
	}

	return resumes, nil
}

func (ur *SqliteResumeRepository) GetResumeByID(resumeID, userID string) (*model.Resume, error) {
	query := `select id, title, created_at, updated_at, is_scheduled from resumes where id = ? and user_id = ?`

	var r model.Resume
	if err := ur.DB.QueryRow(query, resumeID, userID).Scan(
		&r.ID,
		&r.Title,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.IsScheduled,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.Resume{}, err
		}
		return nil, err
	}

	return &r, nil
}

func (ur *SqliteResumeRepository) ToggleScheduling(resumeID, userID string, isScheduled bool) error {
	scheduledValue := 0
	if isScheduled {
		scheduledValue = 1
	}

	query := `update resumes set is_scheduled = ? where id = ? and user_id = ?`

	if _, err := ur.DB.Exec(query, scheduledValue, resumeID, userID); err != nil {
		return err
	}

	return nil

}
