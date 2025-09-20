package model

import "pkg/utils"

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint
	Code         string
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
