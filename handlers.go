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
		log.Println("/home " + err.Error())
		templates.ExecuteTemplate(w, "base", nil)
		return
	}

	resumes, err := getResumesByUserID(db, u)
	if err != nil {
		log.Println("/home " + err.Error())
		templates.ExecuteTemplate(w, "base", nil)
		return
	}
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

func toggleResume(w http.ResponseWriter, r *http.Request) {
	resumeID := r.PathValue("id")
	userID := sessionManager.GetString(r.Context(), "userID")

	var err error

	if err = r.ParseForm(); err != nil {
		http.Error(w, "could not parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	desiredIsScheduled := r.Form.Has("is_scheduled")

	if err = updateResumeScheduling(db, resumeID, userID, desiredIsScheduled); err != nil {
		http.Error(w, "could not update resume scheduling in database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var resume *Resume
	if resume, err = getResumeByID(db, resumeID, userID); err != nil {
		http.Error(w, "could not fetch updated resume from database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "toggle-switch", resume)
}

func updateResumesOnDemand(w http.ResponseWriter, r *http.Request) {
	var hhr, dbr []Resume
	var err error

	userID := sessionManager.GetString(r.Context(), "userID")
	token, err := getTokenByUserID(db, userID)
	if err != nil {
		http.Error(w, "could not get token by user id from database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if hhr, err = HHGetResumes(client, token.AccessToken); err != nil {
		http.Error(w, "could not get resumes from hh: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if dbr, err = getResumesByUserID(db, userID); err != nil {
		http.Error(w, "could not get resumes by user_id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = reconcileResumes(db, hhr, dbr, userID); err != nil {
		http.Error(w, "could not reconcile resumes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
