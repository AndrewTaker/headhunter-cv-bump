package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func home(w http.ResponseWriter, r *http.Request) {
	u := sessionManager.GetString(r.Context(), "userID")
	user, err := getUserByID(db, u)
	if err != nil {
		templates.ExecuteTemplate(w, "base", nil)
		return
	}

	resumes, err := getResumesByUserID(db, u)
	if err != nil {
		templates.ExecuteTemplate(w, "base", nil)
		return
	}
	log.Println(resumes)
	templates.ExecuteTemplate(w, "base", map[string]any{"User": user, "Resumes": resumes})
}

func page(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "base", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	state, err := GenerateState(64)
	if err != nil {
		http.Error(w, "could not generate state string", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "auth_state",
		Value:   state,
		Expires: time.Now().Add(5 * time.Minute),
		Path:    "/",
		// HttpOnly: true, // Optional: prevents client-side script access
		// Secure:   true, // Optional: only send over HTTPS
		// SameSite: http.SameSiteLaxMode, // Optional: prevents CSRF
	})

	redirectURI := fmt.Sprintf("%s://%s:%s/auth/callback", serverHTTP, serverHost, serverPort)

	url := fmt.Sprintf(
		"https://hh.ru/oauth/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s",
		"code", clientID, state, redirectURI,
	)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "bad code", http.StatusBadRequest)
		return
	}

	queryState := r.URL.Query().Get("state")
	cookieState, err := r.Cookie("auth_state")
	if err != nil {
		http.Error(w, "bad state", http.StatusBadRequest)
		return
	}

	if cookieState.Value != queryState {
		http.Error(w, "states dont match", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "auth_state",
		Value:  "",
		MaxAge: -1,
	})

	token, err := HHGetToken(client, code)
	if err != nil {
		http.Error(w, "could not get token "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := HHGetUser(client, token.AccessToken)
	if err != nil {
		http.Error(w, "could not get user "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = createOrUpdateUser(db, user); err != nil {
		http.Error(w, "could not create user to database "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = createOrUpdateTokens(db, *token, code, user.ID); err != nil {
		http.Error(w, "could not create or update tokens to database "+err.Error(), http.StatusBadRequest)
		return
	}

	resumes, err := HHGetResumes(client, token.AccessToken)
	if err != nil {
		http.Error(w, "could not fetch resumes "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = createOrUpdateResumes(db, resumes, user.ID); err != nil {
		http.Error(w, "could not create or update resumes to database "+err.Error(), http.StatusBadRequest)
		return
	}
	sessionManager.Put(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
