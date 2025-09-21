package repository

import (
	"pkg/database"
	"pkg/model"
)

type SchedulerRepository interface {
	GetSchedules() ([]model.JoinedScheduler, error)
	SaveResult(s model.JoinedScheduler, timestamp, err string) error
}

type SqliteSchedulerRepository struct {
	DB *database.DB
}

func NewSqliteSchedulerRepository(db *database.DB) SchedulerRepository {
	return &SqliteSchedulerRepository{DB: db}
}

func (sr *SqliteSchedulerRepository) GetSchedules() ([]model.JoinedScheduler, error) {
	query := `
	select users.id, tokens.access_token, tokens.refresh_token, resumes.id, resumes.title
	from users
	join tokens on users.id = tokens.user_id
	join resumes on users.id = resumes.user_id
	`
	rows, err := sr.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var data []model.JoinedScheduler
	for rows.Next() {
		var s model.JoinedScheduler
		if err := rows.Scan(
			&s.UserID,
			&s.AccessToken,
			&s.RefreshToken,
			&s.ResumeID,
			&s.ResumeTitle,
		); err != nil {
			return nil, err
		}

		data = append(data, s)
	}

	return data, nil
}

func (sr *SqliteSchedulerRepository) SaveResult(s model.JoinedScheduler, timestamp, errors string) error {
	query := `
	insert into scheduler (user_id, resume_id, resume_title, timestamp, error)
	values (?, ?, ?, ?, ?)
	`

	if _, err := sr.DB.Exec(query, s.UserID, s.ResumeID, s.ResumeTitle, timestamp, errors); err != nil {
		return err
	}

	return nil
}
