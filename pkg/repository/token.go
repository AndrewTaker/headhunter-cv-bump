package repository

import "pkg/model"

type TokenRepository interface {
	CreateOrUpdateToken(token model.Token, code, userID string) error
	GetTokenByUserID(userID string) (*model.Token, error)
}
