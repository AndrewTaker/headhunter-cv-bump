package database

import (
	"database/sql"
	"log/slog"
	"net/url"
)

type DB struct {
	*sql.DB
}

func NewSqliteDatabase(path string) (*DB, error) {
	params := url.Values{}
	params.Add("_fk", "on")
	dsn := "file:" + url.PathEscape(path) + "?" + params.Encode()

	db, err := sql.Open("sqlite3", dsn)
	slog.Info(dsn)
	if err != nil {
		return nil, err
	}

	var fk string
	_ = db.QueryRow("pragma foreign_keys").Scan(&fk)
	slog.Info(fk)

	_, err = db.Exec(tables)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

var tables = `
	create table if not exists users (
		id text primary key,
		first_name text,
		last_name text,
		middle_name text
	) without rowid;

	create table if not exists tokens (
		access_token text,
		refresh_token text,
		token_type text,
		expiry integer,
		user_id text unique not null,

		foreign key (user_id) references users(id) on delete cascade
	);

	create table if not exists resumes (
		id text primary key unique,
		alternate_url text,
		title text,
		created_at text,
		updated_at text,
		is_scheduled integer not null default 0,
		user_id text not null,

		foreign key (user_id) references users(id) on delete cascade
	);

	create table if not exists scheduler (
		user_id text,
		resume_id text,
		resume_title text,
		timestamp text,
		error text
	);

	create table if not exists session (
		id text,
		expires_at text,
		user_id text not null,

		foreign key (user_id) references users(id) on delete cascade
	);
`
