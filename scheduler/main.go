package main

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"

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
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("db query err ", err)
	}

	var data []S
	for rows.Next() {
		var s S
		if err := rows.Scan(&s.UserID, &s.AccessToken, &s.ResumeID, &s.ResumeTitle); err != nil {
			s.Error += err.Error()
		}
		s.AccessToken, err = Decrypt(s.AccessToken)
		if err != nil {
			s.Error += err.Error()
		}
		data = append(data, s)
	}

	for _, u := range data {
		log.Println(u.UserID, u.AccessToken)
	}
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
