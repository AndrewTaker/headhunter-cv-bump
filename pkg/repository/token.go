package repository

import (
	"context"
	"database/sql"
	"pkg/database"
	"pkg/model"
	"time"
)

type TokenRepository interface {
	GetTokenByUserID(ctx context.Context, userID string) (*model.Token, error)
	SaveToken(ctx context.Context, token *model.Token, userID string) error
	UpdateToken(ctx context.Context, token *model.Token, userID string) error
}

type SqliteTokenRepository struct {
	DB *database.DB
}

func NewSqliteTokenRepository(db *database.DB) TokenRepository {
	return &SqliteTokenRepository{DB: db}
}

func (tr *SqliteTokenRepository) SaveToken(ctx context.Context, token *model.Token, userID string) error {
	query := `
	insert into tokens (access_token, refresh_token, token_type, expiry, user_id)
	values (?, ?, ?, ?, ?)
	`

	if _, err := tr.DB.Exec(query,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.Expiry.Unix(),
		userID,
	); err != nil {
		return err
	}

	return nil
}

func (tr *SqliteTokenRepository) GetTokenByUserID(ctx context.Context, userID string) (*model.Token, error) {
	query := `select access_token, refresh_token, token_type, expiry from tokens where user_id = ?`

	var t model.Token
	var e int64
	if err := tr.DB.QueryRow(query, userID).Scan(
		&t.AccessToken,
		&t.RefreshToken,
		&t.TokenType,
		&e,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.Token{}, err
		}
		return nil, err
	}

	t.Expiry = time.Unix(e, 0)
	return &t, nil
}

func (tr *SqliteTokenRepository) UpdateToken(ctx context.Context, token *model.Token, userID string) error {
	query := `
	update tokens 
	set access_token = ?, refresh_token = ?, token_type = ?, expiry = ?
	where user_id = ?
	`

	if _, err := tr.DB.Exec(
		query,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.Expiry.Unix(),
		userID,
	); err != nil {
		return err
	}

	return nil
}
