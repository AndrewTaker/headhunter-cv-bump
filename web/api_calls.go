package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type Endpoint struct {
	Me        string
	MyResumes string
}

type HeadHunter struct {
	DB            *DB
	Oauth2Config  *oauth2.Config
	DefaultClient *http.Client
	Endpoint      Endpoint
	Headers       map[string]string
}

func NewHeadHunter(db *DB) *HeadHunter {
	return &HeadHunter{
		DB: db,
		Oauth2Config: &oauth2.Config{
			ClientID:     os.Getenv("HH_CLIENT_ID"),
			ClientSecret: os.Getenv("HH_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("HH_REDIRECT_URL"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://hh.ru/oauth/authorize",
				TokenURL: "https://api.hh.ru/token",
			},
		},
		DefaultClient: &http.Client{Timeout: 15 * time.Second},
		Endpoint: Endpoint{
			Me:        "https://api.hh.ru/me",
			MyResumes: "https://api.hh.ru/resumes/mine",
		},
		Headers: map[string]string{
			"HH-User-Agent": os.Getenv("HH_USER_AGENT"),
		},
	}
}

func (hh *HeadHunter) ApplyHeaders(r *http.Request) {
	for k, v := range hh.Headers {
		r.Header.Set(k, v)
	}
}

// we have time string without offset
// because fuck you that's why
type HHTime time.Time

const timeLayout = "2006-01-02 15:04:05-07:00"

func (hht *HHTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02T15:04:05-0700", s)
	if err != nil {
		return err
	}
	*hht = HHTime(t)
	return nil
}

func (hht HHTime) MarshalJSON() ([]byte, error) {
	t := time.Time(hht)
	return []byte(`"` + t.Format(timeLayout) + `"`), nil
}

func (t HHTime) Format(layout string) string {
	return time.Time(t).Format(time.RFC1123)
}

type HHError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type Resume struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CreatedAt    HHTime `json:"created_at"`
	UpdatedAt    HHTime `json:"updated_at"`
	AlternateURL string `json:"alternate_url"`
	IsScheduled  int
}

func (hh *HeadHunter) GetUser(ctx context.Context, token *oauth2.Token) (*User, error) {

	client := oauth2.NewClient(ctx, token)
	req, err := http.NewRequest("GET", "https://api.hh.ru/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+at)
	hh.ApplyHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status code getUser(): %d %s", resp.StatusCode, bodyBytes)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil

}

func HHGetResumes(client *http.Client, t string) ([]Resume, error) {
	req, err := http.NewRequest("GET", "https://api.hh.ru/resumes/mine", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("HH-User-Agent", "n0thingg@yandex.ru update-cv")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status code getUser(): %d %s", resp.StatusCode, bodyBytes)
	}

	type hhResumesResponse struct {
		Items []Resume `json:"items"`
	}
	var hhr hhResumesResponse
	if err := json.NewDecoder(resp.Body).Decode(&hhr); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return hhr.Items, nil
}
