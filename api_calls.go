package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HHError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint   `json:"expires_in"`
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
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	AlternateURL string `json:"alternate_url"`
}

func HHGetToken(client *http.Client, code string) (*Token, error) {
	req, err := http.NewRequest("POST", "https://api.hh.ru/token", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HH-User-Agent", "n0thingg@yandex.ru update-cv")

	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", "http://localhost:44444/auth/callback")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var hherror HHError
		if err = json.NewDecoder(resp.Body).Decode(&hherror); err != nil {
			return nil, fmt.Errorf("failed to decode nested token response: %w", err)
		}
		if hherror.Error == "invalid_grant" && hherror.ErrorDescription == "code has already been used" {
			return nil, fmt.Errorf("invalud grant, returning")
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status code getToken(): %d %s", resp.StatusCode, bodyBytes)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &token, nil

}

func HHGetUser(client *http.Client, t string) (*User, error) {
	req, err := http.NewRequest("GET", "https://api.hh.ru/me", nil)
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

	// var resumes []Resume
	type hhResumesResponse struct {
		Items []Resume `json:"items"`
	}
	var hhr hhResumesResponse
	if err := json.NewDecoder(resp.Body).Decode(&hhr); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return hhr.Items, nil
}

func HHInvalidateToken(client *http.Client, t string) error {
	req, err := http.NewRequest("DELETE", "https://api.hh.ru/me", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("HH-User-Agent", "n0thingg@yandex.ru update-cv")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status code invalidateToken(): %d %s", resp.StatusCode, bodyBytes)
	}

	return nil
}
