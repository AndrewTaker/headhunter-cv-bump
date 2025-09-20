package repository

import (
	"pkg/database"
	"pkg/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
}

type SqliteUserRepository struct {
	DB *database.DB
}

func NewSqliteUserRepository(db *database.DB) UserRepository {
	return &SqliteUserRepository{DB: db}
}

func (ur *SqliteUserRepository) CreateUser(user *model.User) error          { return nil }
func (ur *SqliteUserRepository) GetUserByID(id string) (*model.User, error) { return nil, nil }
