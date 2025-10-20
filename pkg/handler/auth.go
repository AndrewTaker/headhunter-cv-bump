package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkg/auth"
	"pkg/headhunter"
	"pkg/model"
	"pkg/service"
	"pkg/utils"
	"sync"
	"time"
)

var states = make(map[string]struct{})
var stateMutex sync.Mutex

type AuthHandler struct {
	userService   service.UserService
	tokenService  service.TokenService
	resumeService service.ResumeService
	auth          *auth.AuthRepository
}

func NewAuthHandler(
	us service.UserService,
	ts service.TokenService,
	rs service.ResumeService,
	auth *auth.AuthRepository,
) *AuthHandler {
	return &AuthHandler{us, ts, rs, auth}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	u := model.User{ID: "1", FirstName: "QQQ", LastName: "WWW", MiddleName: "EQWEQE"}
	data, _ := json.Marshal(u)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		http.Error(w, "Could not retrieve cookie", http.StatusBadRequest)
		log.Println("cookie err", err)
		return
	}

	h.auth.InvalidateToken(token.Value)
	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateRandomString(64)

	http.SetCookie(w, &http.Cookie{
		Name:   "sess",
		Value:  "",
		MaxAge: -1,
	})

	if err != nil {
		http.Error(w, "Error generating state string", http.StatusBadRequest)
		return
	}

	stateMutex.Lock()
	states[state] = struct{}{}
	stateMutex.Unlock()

	http.Redirect(w, r, headhunter.GetAuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	queryState := r.URL.Query().Get("state")
	authCode := r.URL.Query().Get("code")

	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		http.Error(w, fmt.Sprintf("Authorization failed: %s", errorMsg), http.StatusForbidden)
		return
	}

	if !validateStateToken(queryState) {
		http.Error(w, "Invalid or missing state parameter. Potential CSRF.", http.StatusForbidden)
		return
	}

	ctx := context.Background()
	token, err := headhunter.ExchangeCodeForToken(r.Context(), authCode)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		http.Error(w, "Could not exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := headhunter.HHOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.hh.ru/me")
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		http.Error(w, "Failed to fetch user data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("HeadHunter API returned non-200 status: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var user headhunter.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode user info JSON: %v", err)
		http.Error(w, "Failed to decode user data.", http.StatusInternalServerError)
		return
	}

	if err := h.userService.CreateOrUpdateUser(&model.User{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, MiddleName: user.MiddleName}); err != nil {
		log.Printf("Failed to save user to database: %v", err)
		http.Error(w, "Failed to save user to database", http.StatusInternalServerError)
		return
	}

	dbToken := model.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	if err := h.tokenService.SaveToken(ctx, &dbToken, user.ID); err != nil {
		log.Printf("Failed to save token to database: %v", err)
		http.Error(w, "Failed to save token to database", http.StatusInternalServerError)
		return
	}

	sessionToken, err := utils.GenerateRandomString(32)
	if err != nil {
		log.Printf("failed to generate random string %v", err)
		http.Error(w, "failed to generate random string", http.StatusInternalServerError)
		return
	}

	h.auth.StoreToken(user.ID, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    sessionToken,
		Expires:  time.Now().Add(5 * time.Minute),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

func validateStateToken(state string) bool {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	if _, ok := states[state]; !ok {
		return false
	}

	delete(states, state)
	return true
}
