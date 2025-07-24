package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	User    *User
	Resumes *[]Resume

	Notification string
	Error        string

	IsLoggedIn bool
}

func home(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		IsLoggedIn: false,
	}

	flash := sessionManager.PopString(r.Context(), "notification")
	if flash != "" {
		data.Notification = flash
	}

	error := sessionManager.PopString(r.Context(), "error")
	if error != "" {
		data.Error = error
	}

	u := sessionManager.GetString(r.Context(), "userID")
	if u != "" {
		user, err := getUserByID(db, u)
		if err != nil {
			log.Printf("ERROR: /home failed to get user %s: %v", u, err)
			data.Error = "Could not load your user profile. Please try logging in again."
		} else {
			data.IsLoggedIn = true
			data.User = user
		}

		if data.User != nil {
			resumes, err := getResumesByUserID(db, u)
			if err != nil {
				log.Printf("ERROR: /home failed to get resumes for user %s: %v", u, err)
				if data.Error == "" {
					data.Error = "Could not load your resumes. Please try refreshing."
				} else {
					data.Error += " Also, could not load your resumes."
				}
			} else {
				data.Resumes = &resumes
			}
		} else {
			data.Resumes = nil
		}
	} else {
		data.IsLoggedIn = false
		data.User = nil
		data.Resumes = nil
	}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("ERROR: /home failed to execute template: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}
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
		Path:   "/",
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
	var errMsg string

	if err = r.ParseForm(); err != nil {
		log.Println(err)
		errMsg += " Error reading input. Try again."
	}

	desiredIsScheduled := r.Form.Has("is_scheduled")

	if err = updateResumeScheduling(db, resumeID, userID, desiredIsScheduled); err != nil {
		log.Println(err)
		errMsg += " Could not update. Try again."
	}

	var resume *Resume
	if resume, err = getResumeByID(db, resumeID, userID); err != nil {
		log.Println(err)
		errMsg += " Could not update. Try again."
	}

	sessionManager.Put(r.Context(), "error", errMsg)
	templates.ExecuteTemplate(w, "toggle-switch", resume)
}

func updateResumesOnDemand(w http.ResponseWriter, r *http.Request) {
	var hhr, dbr []Resume
	var err error
	var errMsg string

	userID := sessionManager.GetString(r.Context(), "userID")
	if userID == "" {
		sessionManager.Put(r.Context(), "error", "Not logged in.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := getTokenByUserID(db, userID)
	if err != nil {
		log.Println("getTokenByUserID ", err)
		errMsg += " Could not identify. Try again."
	}

	if hhr, err = HHGetResumes(client, token.AccessToken); err != nil {
		log.Println("HHGetResumes ", err)
		errMsg += " Could not get resumes from hh api. Try again."
	}

	if dbr, err = getResumesByUserID(db, userID); err != nil {
		log.Println("getResumesByUserID", err)
		errMsg += " Could not get user. Try again."
	}

	if err = reconcileResumes(db, hhr, dbr, userID); err != nil {
		log.Println("reconcileResumes ", err)
		errMsg += " Error deleting data. Try again."
	}

	sessionManager.Put(r.Context(), "error", errMsg)
	sessionManager.Put(r.Context(), "notification", "Updated")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func openModal(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "modal", nil)
}

func closeModal(w http.ResponseWriter, r *http.Request) {
	return
}

func invalidateUserData(w http.ResponseWriter, r *http.Request) {
	userID := sessionManager.GetString(r.Context(), "userID")

	if userID == "" {
		sessionManager.Put(r.Context(), "error", "Not logged in.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var err error
	var t *Token
	var errMsg string
	if t, err = getTokenByUserID(db, userID); err != nil {
		log.Println(err)
		errMsg += " Could not get credentials. Try again."
	}

	if err = HHInvalidateToken(client, t.AccessToken); err != nil {
		log.Println(err)
		errMsg += " Could not invalidate data from headhunter api. Contact to invalidate manually or try again."
	}

	if err = deleteUserByID(db, userID); err != nil {
		log.Println(err)
		errMsg += " Could not delete user data. Try again."
	}

	sessionManager.Remove(r.Context(), "userID")
	sessionManager.Put(r.Context(), "notification", "your data was deleted.")
	sessionManager.Put(r.Context(), "error", errMsg)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)

}

func logout(w http.ResponseWriter, r *http.Request) {
	sessionManager.Destroy(r.Context())
	sessionManager.RenewToken(r.Context())
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
