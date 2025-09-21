package repository

import (
	"database/sql"
	"pkg/database"
	"pkg/model"
)

type TokenRepository interface {
	CreateOrUpdateToken(token *model.Token, code, userID string) error
	GetTokenByUserID(userID string) (*model.Token, error)
	UpdateToken(accessToken, refreshToken, userID string) error
}

type SqliteTokenRepository struct {
	DB *database.DB
}

func NewSqliteTokenRepository(db *database.DB) TokenRepository {
	return &SqliteTokenRepository{DB: db}
}

func (tr *SqliteTokenRepository) CreateOrUpdateToken(token *model.Token, code, userID string) error {
	query := `
	insert into tokens (access_token, refresh_token, expires_in, code, user_id) values (?, ?, ?, ?, ?)
	on conflict(user_id) do update set
	access_token = excluded.access_token,
	refresh_token = excluded.refresh_token,
	expires_in = excluded.expires_in,
	code = excluded.code
	`

	if _, err := tr.DB.Exec(query,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresIn,
		code,
		userID,
	); err != nil {
		return err
	}

	return nil
}

func (tr *SqliteTokenRepository) GetTokenByUserID(userID string) (*model.Token, error) {
	query := `select access_token, refresh_token, expires_in from tokens where user_id = ?`

	var t model.Token
	if err := tr.DB.QueryRow(query, userID).Scan(
		&t.AccessToken,
		&t.RefreshToken,
		&t.ExpiresIn,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.Token{}, err
		}
		return nil, err
	}

	return &t, nil
}

func (tr *SqliteTokenRepository) UpdateToken(accessToken, refreshToken, userID string) error {
	query := `
	update tokens 
	set access_token = ?, refresh_token = ?
	where user_id = ?
	`

	if _, err := tr.DB.Exec(query, accessToken, refreshToken, userID); err != nil {
		return err
	}

	return nil
}
