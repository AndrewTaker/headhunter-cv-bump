package main

import (
	"database/sql"
	"fmt"
)

func db_init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./hh.db")
	if err != nil {
		return nil, err
	}

	tables := `
	create table if not exists users (
		id text primary key,
		first_name text,
		last_name text,
		middle_name text
	);

	create table if not exists tokens (
		access_token text,
		refresh_token text,
		expires_in integer,
		code text unique,
		user_id text unique,
		
		foreign key (user_id) references users(id) on delete cascade
	);

	create table if not exists resumes (
		id text primary key unique,
		alternate_url text,
		title text,
		created_at text,
		updated_at text,
		user_id text,
		is_scheduled integer not null default 0,

		foreign key (user_id) references users(id) on delete cascade
	);
	`

	_, err = db.Exec(tables)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createOrUpdateUser(db *sql.DB, user *User) error {
	query := `
	insert into users (id, first_name, last_name, middle_name) values (?, ?, ?, ?)
	on conflict(id) do update set
	first_name = excluded.first_name,
	last_name = excluded.last_name,
	middle_name = excluded.middle_name
	`
	_, err := db.Exec(query, user.ID, user.FirstName, user.LastName, user.MiddleName)
	if err != nil {
		return err
	}

	return nil
}

func getUserByID(db *sql.DB, userID string) (*User, error) {
	query := `select id, first_name, last_name, middle_name from users where id = ?`
	var u User
	if err := db.QueryRow(query, userID).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.MiddleName,
	); err != nil {
		return nil, err
	}

	return &u, nil
}

func getTokenByCode(db *sql.DB, code string) (*Token, error) {
	query := `select access_token, refresh_token, expires_in from tokens where code = ?`
	var t Token
	if err := db.QueryRow(query, code).Scan(
		&t.AccessToken,
		&t.RefreshToken,
		&t.ExpiresIn,
	); err != nil {
		return nil, err
	}

	return &t, nil
}

func createOrUpdateTokens(db *sql.DB, tokens Token, code string, userID string) error {
	query := `
	insert into tokens (access_token, refresh_token, expires_in, code, user_id) values (?, ?, ?, ?, ?)
	on conflict(user_id) do update set
	access_token = excluded.access_token,
	refresh_token = excluded.refresh_token,
	expires_in = excluded.expires_in,
	code = excluded.code
	`
	_, err := db.Exec(query, tokens.AccessToken, tokens.RefreshToken, tokens.ExpiresIn, code, userID)
	if err != nil {
		return err
	}

	return nil
}

func createOrUpdateResumes(db *sql.DB, resumes []Resume, userID string) error {
	query := `
	insert into resumes (id, title, created_at, updated_at, user_id) values (?, ?, ?, ?, ?)
	on conflict(id) do update set
	title = excluded.title,
	created_at = excluded.created_at,
	updated_at = excluded.updated_at,
	user_id = excluded.user_id;
	`
	tx, err := db.Begin()
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

func getResumesByUserID(db *sql.DB, userID string) ([]Resume, error) {
	query := "select id, title, created_at, updated_at from resumes where user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	var resumes []Resume
	for rows.Next() {
		var r Resume
		if err := rows.Scan(&r.ID, &r.Title, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		resumes = append(resumes, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resumes, nil

}
