package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrInternal      = errors.New("internal error occured")
)

type ErrorResponse struct {
	E string `json:"error"`
}

type ResumeResponseSingle struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	AlternateUrl string    `json:"alternate_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsScheduled  int       `json:"is_scheduled"`
}

type ResumeResponseMany struct {
	Resumes []ResumeResponseSingle `json:"resumes"`
}

type UserResponse struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

func JsonResponseErr(w http.ResponseWriter, r *http.Request, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := ErrorResponse{message}
	b, _ := json.Marshal(e)
	w.Write(b)
}

func JsonResponse(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	b, _ := json.Marshal(data)
	w.Write(b)
}
