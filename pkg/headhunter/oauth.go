package headhunter

import (
	"context"
	"os"

	"golang.org/x/oauth2"
)

var HHOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("HH_CLIENT_ID"),
	ClientSecret: os.Getenv("HH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("HH_REDIRECT_URI"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://hh.ru/oauth/authorize",
		TokenURL: "https://hh.ru/oauth/token",
	},
}

func GetAuthCodeURL() string {
	return HHOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return HHOauthConfig.Exchange(ctx, code)
}
