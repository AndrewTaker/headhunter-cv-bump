package main

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal("db err ", err)
	}

	query := `
	select users.id, tokens.access_token, resumes.id, resumes.title
	from users
	join tokens on users.id = tokens.user_id
	join resumes on users.id = resumes.user_id;
	`

	type S struct {
		UserID      string
		AccessToken string
		ResumeID    string
		ResumeTitle string
		Error       string
		Timestamp   string
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("db query err ", err)
	}

	var data []S
	for rows.Next() {
		var s S
		var at string
		if err := rows.Scan(&s.UserID, &at, &s.ResumeID, &s.ResumeTitle); err != nil {
			s.Error += err.Error()
		}

		s.AccessToken, err = Decrypt(at)
		if err != nil {
			s.Error += err.Error()
		}

		data = append(data, s)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	for _, u := range data {
		timestamp, err := bump(client, u.AccessToken, u.ResumeID)
		if err != nil {
			u.Error = err.Error()
		}
		u.Timestamp = timestamp

		_, err = db.Exec(
			`insert into scheduler user_id, resume_id, resume_title, timestamp, error`,
			u.UserID, u.ResumeID, u.ResumeTitle, u.Timestamp, u.Error,
		)
		if err != nil {
			log.Println("err saving result" + err.Error())
		}
	}
}

func bump(client *http.Client, at, rid string) (string, error) {
	url := fmt.Sprintf("https://api.hh.ru/resumes/%s/publish", rid)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("HH-User-Agent", os.Getenv("HH_USER_AGENT"))
	req.Header.Add("Authorization", "Bearer "+at)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("non 204 returned" + string(body))
	}

	return time.Now().Format(time.RFC3339), nil
}

func Decrypt(encryptedString string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", fmt.Errorf("could not decode base64: %w", err)
	}

	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		return "", fmt.Errorf("could not create cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("could not create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short (missing nonce)")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("could not decrypt: %w", err)
	}

	return string(plaintext), nil
}
