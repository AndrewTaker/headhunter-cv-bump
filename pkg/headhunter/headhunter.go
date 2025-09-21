package headhunter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type HHTransport struct {
	headers http.Header
	base    http.RoundTripper
}

func (hht *HHTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())

	for key, values := range hht.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return hht.base.RoundTrip(req)
}

type HHClient struct {
	baseURL string
	client  *http.Client
	config  *oauth2.Config
}

func NewHHClient(ctx context.Context, token *oauth2.Token) *HHClient {
	config := HHOauthConfig
	tokenSource := config.TokenSource(ctx, token)
	httpClient := oauth2.NewClient(ctx, tokenSource)

	// prepend headers via roundtrip
	// this one is mandatory by hh documentation
	headers := make(http.Header)
	headers.Set("HH-User-Agent", os.Getenv("HH_USER_AGENT"))

	httpClient.Transport = &HHTransport{
		headers: headers,
		base:    httpClient.Transport,
	}

	return &HHClient{
		client:  httpClient,
		baseURL: "https://api.hh.ru",
		config:  config,
	}
}

func (hh *HHClient) GetUser(ctx context.Context) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, hh.baseURL+"/me", nil)
	if err != nil {
		return nil, err
	}

	resp, err := hh.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hh.GetUser: bad status code %d", resp.StatusCode)
	}

	var user User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
