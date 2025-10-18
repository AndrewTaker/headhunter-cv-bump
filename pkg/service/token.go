package service

import (
	"context"
	"fmt"
	"log"
	"pkg/model"
	"pkg/repository"

	"golang.org/x/oauth2"
)

type TokenService interface {
	SaveToken(ctx context.Context, token *model.Token, userID string) error
	UpdateToken(ctx context.Context, token *model.Token, userID string) error
	GetTokenByUserID(ctx context.Context, userID string) (*model.Token, error)
}

type TokenServiceImpl struct {
	tokenRepo repository.TokenRepository
}

func NewTokenService(tr repository.TokenRepository) TokenService {
	return &TokenServiceImpl{tokenRepo: tr}
}

func (ts *TokenServiceImpl) SaveToken(ctx context.Context, token *model.Token, userID string) error {
	if err := token.Encrypt(); err != nil {
		return err
	}

	return ts.tokenRepo.SaveToken(ctx, token, userID)
}

func (ts *TokenServiceImpl) UpdateToken(ctx context.Context, token *model.Token, userID string) error {
	if err := token.Encrypt(); err != nil {
		return err
	}

	return ts.tokenRepo.UpdateToken(ctx, token, userID)
}

func (ts *TokenServiceImpl) GetTokenByUserID(ctx context.Context, userID string) (*model.Token, error) {
	token, err := ts.tokenRepo.GetTokenByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := token.Decrypt(); err != nil {
		return nil, err
	}

	return token, nil
}

type TokenSourceSqlite struct {
	userID     string
	repository repository.TokenRepository
	config     *oauth2.Config
}

func NewTokenSourceSqlite(userID string, repository repository.TokenRepository, config *oauth2.Config) oauth2.TokenSource {
	return &TokenSourceSqlite{userID, repository, config}
}

func (tss *TokenSourceSqlite) Token() (*oauth2.Token, error) {
	ctx := context.Background()

	dbToken, err := tss.repository.GetTokenByUserID(ctx, tss.userID)
	if err != nil {
		return nil, fmt.Errorf("could not load token from database %v", err)
	}

	oauth2Token := dbToken.ToOauth2Token()
	if oauth2Token.Valid() {
		return oauth2Token, nil
	}

	tokenRefresher := tss.config.TokenSource(ctx, oauth2Token)
	newToken, err := tokenRefresher.Token()
	if err != nil {
		return nil, fmt.Errorf("could not refresh token %v", err)
	}

	dbToken.AccessToken = newToken.AccessToken
	dbToken.RefreshToken = newToken.RefreshToken
	dbToken.TokenType = newToken.TokenType
	dbToken.RefreshToken = newToken.RefreshToken

	if err := tss.repository.UpdateToken(ctx, dbToken, tss.userID); err != nil {
		log.Printf("although token was renewed, it was not saved %v", err)
	}

	return newToken, nil
}
