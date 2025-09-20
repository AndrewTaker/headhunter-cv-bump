package model

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint
	Code         string
	TokenType    string
}
