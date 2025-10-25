package headhunter

import (
	"net/url"
)

func GenerateAuthUrl(state, clientID string) string {
	base := "https://hh.ru/oauth/authorize"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("state", state)

	u, err := url.Parse(base)
	if err != nil {
		return ""
	}

	u.RawQuery = params.Encode()

	return u.String()
}
