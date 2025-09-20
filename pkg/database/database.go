package database

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

func NewSqliteDatabase(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(tables)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
