package repository

import (
	"database/sql"
	"pkg/database"
	"pkg/model"
)

type UserRepository interface {
	CreateOrUpdateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	DeleteUserByID(id string) error
}

type SqliteUserRepository struct {
	DB *database.DB
}

func NewSqliteUserRepository(db *database.DB) UserRepository {
	return &SqliteUserRepository{DB: db}
}

func (ur *SqliteUserRepository) CreateOrUpdateUser(user *model.User) error {
	query := `
	insert into users (id, first_name, last_name, middle_name) values (?, ?, ?, ?)
	on conflict(id) do update set
	first_name = excluded.first_name,
	last_name = excluded.last_name,
	middle_name = excluded.middle_name
	`

	if _, err := ur.DB.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.MiddleName,
	); err != nil {
		return err
	}

	return nil
}

func (ur *SqliteUserRepository) GetUserByID(id string) (*model.User, error) {
	query := `select id, first_name, last_name, middle_name from users where id = ?`

	var u model.User
	if err := ur.DB.QueryRow(query, id).Scan(
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

func (ur *SqliteUserRepository) DeleteUserByID(id string) error {
	query := `delete from users where id = ?`

	if _, err := ur.DB.Exec(query, id); err != nil {
		return err
	}

	return nil
}
