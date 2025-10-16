package auth

import (
	"log"
	"sync"
	"time"
)

var mutex sync.RWMutex

type SessionData struct {
	UserID    string
	ExpiresAt time.Time
}

type AuthRepository struct {
	Session map[string]SessionData
}

const TokenTTL = 5 * time.Minute
const CleanupInterval = 1 * time.Minute

func NewAuthRepository() *AuthRepository {
	ar := &AuthRepository{make(map[string]SessionData)}
	go ar.StartCleanupRoutine()
	return ar
}

func (ar *AuthRepository) StoreToken(userID, token string) {
	mutex.Lock()
	defer mutex.Unlock()

	ar.Session[token] = SessionData{userID, time.Now().Add(TokenTTL)}
}

func (ar *AuthRepository) InvalidateToken(token string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(ar.Session, token)
}

func (ar *AuthRepository) GetUserByToken(token string) string {
	mutex.RLock()
	data, exists := ar.Session[token]
	mutex.RUnlock()

	if exists && time.Now().Before(data.ExpiresAt) {
		return data.UserID
	}
	return ""
}

func (ar *AuthRepository) IsPresent(token string) bool {
	mutex.RLock()
	data, exists := ar.Session[token]
	mutex.RUnlock()

	return exists && time.Now().Before(data.ExpiresAt)
}

func (ar *AuthRepository) StartCleanupRoutine() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	log.Println("Cleanup routine started.")

	for range ticker.C {
		ar.cleanupExpiredTokens()
	}
}

func (ar *AuthRepository) cleanupExpiredTokens() {
	now := time.Now()
	log.Printf("[%s] [INFO] AuthRepository: Running cleanup job.", now.Format(time.RFC3339))
	mutex.Lock()
	defer mutex.Unlock()

	for token, data := range ar.Session {
		if now.After(data.ExpiresAt) {
			delete(ar.Session, token)
		}
	}
}
