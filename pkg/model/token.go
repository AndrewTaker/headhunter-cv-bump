package model

import (
	"pkg/utils"
	"time"

	"golang.org/x/oauth2"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	TokenType    string
}

func (t *Token) Encrypt() error {
	var err error

	if t.AccessToken, err = utils.Encrypt(t.AccessToken); err != nil {
		return err
	}

	if t.RefreshToken, err = utils.Encrypt(t.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (t *Token) Decrypt() error {
	var err error

	if t.AccessToken, err = utils.Decrypt(t.AccessToken); err != nil {
		return err
	}

	if t.RefreshToken, err = utils.Decrypt(t.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (t *Token) ToOauth2Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenType:    t.TokenType,
		Expiry:       t.Expiry,
	}
}
