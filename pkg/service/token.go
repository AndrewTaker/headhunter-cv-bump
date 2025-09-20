package service

import (
	"pkg/model"
	"pkg/repository"
)

type TokenService interface {
	CreateOrUpdateToken(token *model.Token, code, userID string) error
	GetToken(userID string) (*model.Token, error)
}

type TokenServiceImpl struct {
	tokenRepo repository.TokenRepository
}

func NewTokenService(tr repository.TokenRepository) TokenService {
	return &TokenServiceImpl{tokenRepo: tr}
}

func (ts *TokenServiceImpl) CreateOrUpdateToken(token *model.Token, code, userID string) error {
	if err := token.Encrypt(); err != nil {
		return err
	}

	return ts.tokenRepo.CreateOrUpdateToken(token, code, userID)
}

func (ts *TokenServiceImpl) GetToken(userID string) (*model.Token, error) {
	token, err := ts.tokenRepo.GetTokenByUserID(userID)
	if err != nil {
		return nil, err
	}

	if err := token.Decrypt(); err != nil {
		return nil, err
	}

	return token, nil
}
