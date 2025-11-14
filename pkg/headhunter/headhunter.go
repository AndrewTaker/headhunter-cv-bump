package headhunter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	ErrUnauthorized   = errors.New("Unauthorized")
	ErrHHTokenExpired = errors.New("token_expired")
)

type HHClient struct {
	baseURL string
	client  *http.Client
	AT      *string
	RT      *string

	log *slog.Logger
}

func NewHHClient(ctx context.Context, at, rt *string) *HHClient {
	logger := slog.Default().With(slog.String("log_type", "hhclient"))
	return &HHClient{
		client:  &http.Client{Timeout: time.Second * 15},
		baseURL: "https://api.hh.ru",
		AT:      at,
		RT:      rt,
		log:     logger,
	}
}

func (hh *HHClient) GetUserResumes(ctx context.Context) ([]Resume, error) {
	const path = "/resumes/mine"
	const method = http.MethodGet

	log := hh.log.With(
		slog.String("target_path", path),
		slog.String("method", http.MethodGet),
	)

	req, err := http.NewRequestWithContext(ctx, method, hh.baseURL+path, nil)
	if err != nil {
		log.Error("failed request creation", "error", err)
		return nil, err
	}

	req.Header.Set("HH-User-Agent", os.Getenv("HH_USER_AGENT"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *hh.AT))

	resp, err := hh.client.Do(req)
	if err != nil {
		log.Error("failed response", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("bad status code", slog.String("body", string(body)))
		return nil, fmt.Errorf("hh.GetUserResumes: bad status code %d", resp.StatusCode)
	}

	var hhr struct {
		Items []Resume `json:"items"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&hhr); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Error("json decoding error",
			"error", err,
			slog.String("body", string(body)),
		)
		return nil, err
	}

	return hhr.Items, nil
}

func (hh *HHClient) GetUser(ctx context.Context) (*User, error) {
	const path = "/me"
	const method = http.MethodGet

	log := hh.log.With(
		slog.String("target_path", path),
		slog.String("method", http.MethodGet),
	)

	req, err := http.NewRequestWithContext(ctx, method, hh.baseURL+path, nil)
	if err != nil {
		log.Error("failed request creation", "error", err)
		return nil, err
	}

	req.Header.Set("HH-User-Agent", os.Getenv("HH_USER_AGENT"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *hh.AT))

	resp, err := hh.client.Do(req)
	if err != nil {
		log.Error("failed response", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		var e AuthError
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Error("unauth error", "decoded err", err)
			return nil, ErrUnauthorized
		}

		for _, hhe := range e.Errors {
			if *hhe.Value == "token_expired" {
				return nil, ErrHHTokenExpired
			}
		}
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Error("bad status code", slog.String("body", string(body)))
		return nil, fmt.Errorf("hh.GetUser: bad status code %d", resp.StatusCode)
	}

	var user User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Error("json decoding error", "error", err)
		return nil, err
	}

	return &user, nil
}

func (hh *HHClient) AuthExachangeCodeForToken(ctx context.Context, code string) (*Token, error) {
	const path = "/token"
	const method = http.MethodPost

	log := hh.log.With(
		slog.String("target_path", path),
		slog.String("method", http.MethodGet),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hh.baseURL+"/token", nil)
	if err != nil {
		log.Error("failed request creation", "error", err)
		return nil, err
	}

	req.Header.Set("HH-User-Agent", os.Getenv("HH_USER_AGENT"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("client_id", os.Getenv("HH_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("HH_CLIENT_SECRET"))
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	resp, err := hh.client.Do(req)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Error("failed response",
			"error", err,
			slog.String("body", string(body)),
		)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("bad status code",
			"error", err,
			slog.String("body", string(body)),
		)
		return nil, fmt.Errorf("hh.AuthExchangeCodeForToken: bad status code %d %s", resp.StatusCode, string(body))
	}

	var t Token
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		log.Error("json decoding error", "error", err)
		return nil, err
	}

	return &t, nil
}

func (hh *HHClient) RefreshToken(ctx context.Context) (*Token, error) {
	const path = "/token"
	const method = http.MethodPost

	log := hh.log.With(
		slog.String("target_path", path),
		slog.String("method", method),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hh.baseURL+path, nil)
	if err != nil {
		log.Error("failed request creation", "error", err)
		return nil, err
	}

	req.Header.Set("HH-User-Agent", os.Getenv("HH_USER_AGENT"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("refresh_token", *hh.RT)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()

	resp, err := hh.client.Do(req)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Error("failed response",
			"error", err,
			slog.String("body", string(body)),
		)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("bad status code",
			"error", err,
			slog.String("body", string(body)),
		)
		return nil, fmt.Errorf("hh.RefreshToken: bad status code %d %s", resp.StatusCode, string(body))
	}

	var t Token
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		log.Error("json decoding error", "error", err)
		return nil, err
	}

	return &t, nil
}

func (hh *HHClient) BumpResume(ctx context.Context, id string) error {
	var path string = fmt.Sprintf("/resumes/%s/publish", id)
	const method = http.MethodPost

	log := hh.log.With(
		slog.String("target_path", path),
		slog.String("method", method),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hh.baseURL+path, nil)
	if err != nil {
		log.Error("failed request creation", "error", err)
		return err
	}

	req.Header.Set("HH-User-Agent", os.Getenv("HH_USER_AGENT"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *hh.AT))

	resp, err := hh.client.Do(req)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Error("failed response",
			"error", err,
			slog.String("body", string(body)),
		)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		var e AuthError
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Error("unauth error", "decoded err", err)
			return ErrUnauthorized
		}

		for _, hhe := range e.Errors {
			if *hhe.Value == "token_expired" {
				return ErrHHTokenExpired
			}
		}
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("bad status code",
			"error", err,
			slog.String("body", string(body)),
		)
		return fmt.Errorf("hh.BumpResume: bad status code %d %v", resp.StatusCode, err)
	}

	return nil
}
