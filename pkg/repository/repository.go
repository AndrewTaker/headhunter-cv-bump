package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"pkg/database"
	"pkg/model"
	"time"
)

type SqliteRepository struct {
	DB  *database.DB
	log *slog.Logger
}

func NewSqliteRepository(db *database.DB) *SqliteRepository {
	logger := slog.Default().With(slog.String("log_type", "database"))
	return &SqliteRepository{DB: db, log: logger}
}

func (sr *SqliteRepository) ResumeCreateOrUpdateBatch(resumes []model.Resume, userID string) error {
	const function = "ResumeCreateOrUpdateBatch"
	log := sr.log.With("function", function)

	query := `
	insert into resumes (id, title, alternate_url, created_at, updated_at, user_id) values (?, ?, ?, ?, ?, ?)
	on conflict(id) do update set
	title = excluded.title,
	alternate_url = excluded.alternate_url,
	created_at = excluded.created_at,
	updated_at = excluded.updated_at,
	user_id = excluded.user_id;
	`
	tx, err := sr.DB.Begin()
	if err != nil {
		log.Error("transaction error", "error", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error("transaction error", "error", err)
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
			log.Error("qquery exec error", "error", err)
			return fmt.Errorf("failed to execute statement for resume ID %s: %w", resume.ID, err)
		}
	}
	tx.Commit()

	return nil
}

func (sr *SqliteRepository) ResumeGetByUserIDBatch(userID string) ([]model.Resume, error) {
	const function = "ResumeGetByUserIDBatch"
	log := sr.log.With("function", function)

	query := "select id, title, alternate_url, created_at, updated_at, is_scheduled from resumes where user_id = ?"

	rows, err := sr.DB.Query(query, userID)
	if err != nil {
		log.Error("query error", "error", err)
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
			log.Error("scan error", "error", err)
			return nil, err
		}
		resumes = append(resumes, r)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows error", "error", err)
		return []model.Resume{}, err
	}

	if err == sql.ErrNoRows {
		return []model.Resume{}, nil
	}

	return resumes, nil
}

func (sr *SqliteRepository) ResumeGetByID(resumeID, userID string) (*model.Resume, error) {
	query := `select id, title, created_at, updated_at, is_scheduled from resumes where id = ? and user_id = ?`

	var r model.Resume
	if err := sr.DB.QueryRow(query, resumeID, userID).Scan(
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

func (sr *SqliteRepository) ResumeToggleScheduling(resumeID, userID string, isScheduled bool) error {
	scheduledValue := 0
	if isScheduled {
		scheduledValue = 1
	}

	query := `update resumes set is_scheduled = ? where id = ? and user_id = ?`

	if _, err := sr.DB.Exec(query, scheduledValue, resumeID, userID); err != nil {
		return err
	}

	return nil

}

func (sr *SqliteRepository) ResumeDeleteByUserID(resumes []model.Resume, userID string) error {
	tx, err := sr.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`delete from resumes where id = ? and user_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, resume := range resumes {
		_, err := stmt.Exec(resume.ID, userID)
		if err != nil {
			return fmt.Errorf("failed to execute statement for resume ID %s: %w", resume.ID, err)
		}
	}
	tx.Commit()

	return nil
}

func (sr *SqliteRepository) UserCreateOrUpdate(user *model.User) error {
	query := `
	insert into users (id, first_name, last_name, middle_name) values (?, ?, ?, ?)
	on conflict(id) do update set
	first_name = excluded.first_name,
	last_name = excluded.last_name,
	middle_name = excluded.middle_name
	`

	if _, err := sr.DB.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.MiddleName,
	); err != nil {
		return err
	}

	return nil
}

func (sr *SqliteRepository) UserGetBySessionID(ctx context.Context, sessID string) (*model.User, error) {
	const function = "UserGetBySessionID"
	log := sr.log.With("function", function)

	query := `
	select id, first_name, last_name, middle_name from users
	where id = (select user_id from session where id = ?)
	`

	var u model.User
	if err := sr.DB.QueryRow(query, sessID).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.MiddleName,
	); err != nil {
		log.Error("scan error", "error", err)
		return nil, err
	}

	return &u, nil
}
func (sr *SqliteRepository) UserGetByID(id string) (*model.User, error) {
	query := `select id, first_name, last_name, middle_name from users where id = ?`

	var u model.User
	if err := sr.DB.QueryRow(query, id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.MiddleName,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.User{}, nil
		}
		return nil, err
	}

	return &u, nil
}

func (sr *SqliteRepository) UserDeleteByID(id string) error {
	query := `delete from users where id = ?`

	if _, err := sr.DB.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (sr *SqliteRepository) ScheduleGetBatch() ([]model.JoinedScheduler, error) {
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

func (sr *SqliteRepository) ScheduleSave(s model.JoinedScheduler, timestamp, errors string) error {
	query := `
	insert into scheduler (user_id, resume_id, resume_title, timestamp, error)
	values (?, ?, ?, ?, ?)
	`

	if _, err := sr.DB.Exec(query, s.UserID, s.ResumeID, s.ResumeTitle, timestamp, errors); err != nil {
		return err
	}

	return nil
}

func (tr *SqliteRepository) TokenSaveOrCreate(ctx context.Context, token *model.Token, userID string) error {
	query := `
	insert into tokens (access_token, refresh_token, user_id) values (?, ?, ?)
	on conflict(user_id) do update set
	access_token = excluded.access_token,
	refresh_token = excluded.refresh_token
	`

	if _, err := tr.DB.Exec(query,
		token.AccessToken,
		token.RefreshToken,
		userID,
	); err != nil {
		return err
	}

	return nil
}

func (sr *SqliteRepository) TokenGetByUserID(ctx context.Context, userID string) (*model.Token, error) {
	query := `select access_token, refresh_token from tokens where user_id = ?`

	var t model.Token
	if err := sr.DB.QueryRow(query, userID).Scan(
		&t.AccessToken,
		&t.RefreshToken,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &t, nil
}

func (sr *SqliteRepository) SessionSave(ctx context.Context, sessID, userID string, expiresAt time.Time) error {
	query := `
	insert into session (id, expires_at, user_id) values (?, ?, ?)
	`

	if _, err := sr.DB.Exec(query, sessID, expiresAt, userID); err != nil {
		return err
	}

	return nil
}

func (sr *SqliteRepository) SessionDelete(ctx context.Context, sessID string) error {
	query := `delete from session where id = ?`

	if _, err := sr.DB.Exec(query, sessID); err != nil {
		return err
	}

	return nil
}
