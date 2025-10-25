package handler

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
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
	service *service.SqliteService
}

func NewAuthHandler(s *service.SqliteService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	sess, err := r.Cookie("sess")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	h.service.DeleteSession(r.Context(), sess.Value)
	http.SetCookie(w, &http.Cookie{
		Name:   "sess",
		Value:  "",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusNoContent)
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

	u := headhunter.GenerateAuthUrl(state, os.Getenv("HH_CLIENT_ID"))

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
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
	client := headhunter.NewHHClient(ctx, nil, nil)
	token, err := client.AuthExachangeCodeForToken(r.Context(), authCode)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		http.Error(w, "Could not exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client.AT = &token.AccessToken
	client.RT = &token.RefreshToken

	user, err := client.GetUser(ctx)
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		http.Error(w, "Failed to fetch user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.CreateOrUpdateUser(&model.User{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
	}); err != nil {
		log.Printf("Failed to save user to database: %v", err)
		http.Error(w, "Failed to save user to database", http.StatusInternalServerError)
		return
	}

	if err := h.service.CreateOrUpdateToken(
		ctx,
		&model.Token{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken},
		user.ID,
	); err != nil {
		log.Printf("failed to generate random string %v", err)
		http.Error(w, "failed to generate random string", http.StatusInternalServerError)
		return
	}

	resumes, err := client.GetUserResumes(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	var dbResumes []model.Resume
	for _, resume := range resumes {
		dbResumes = append(dbResumes, model.Resume{
			ID:           resume.ID,
			AlternateURL: resume.AlternateURL,
			Title:        resume.Title,
			CreatedAt:    model.HHTime(resume.CreatedAt),
			UpdatedAt:    model.HHTime(resume.UpdatedAt),
		})
	}
	if err := h.service.CreateOrUpdateResumes(dbResumes, user.ID); err != nil {
		log.Println(err)
		return
	}

	sessionToken, _ := utils.GenerateRandomString(32)
	expiresAt := time.Now().Add(10 * time.Minute)

	if err := h.service.SaveSession(ctx, sessionToken, user.ID, expiresAt); err != nil {
		log.Printf("failed to save session to db %v", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    sessionToken,
		Expires:  expiresAt,
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
