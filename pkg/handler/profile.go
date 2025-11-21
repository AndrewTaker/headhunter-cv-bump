package handler

import (
	"log/slog"
	"net/http"
	"pkg/service"
	"time"

	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	service *service.SqliteService
	log     *slog.Logger
}

func NewProfileHandler(s *service.SqliteService) *ProfileHandler {
	logger := slog.Default().With(slog.String("log_type", "hhclient"))
	return &ProfileHandler{service: s, log: logger}
}

func (h *ProfileHandler) DeleteUserData(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		JsonResponseErr(w, r, http.StatusForbidden, ErrNotAuthorized.Error())
		return
	}

	user, err := h.service.GetUserBySession(r.Context(), token.Value)
	if err != nil {
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	if err := h.service.DeleteUser(user.ID); err != nil {
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	JsonResponse(w, r, http.StatusOK, struct{}{})
}

func (h *ProfileHandler) Me(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		JsonResponseErr(w, r, http.StatusForbidden, ErrNotAuthorized.Error())
		return
	}

	user, err := h.service.GetUserBySession(r.Context(), token.Value)
	if err != nil {
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	hr := UserResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
	}

	JsonResponse(w, r, http.StatusOK, hr)

}

func (h *ProfileHandler) ToggleResume(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		slog.Error(err.Error())
		JsonResponseErr(w, r, http.StatusForbidden, ErrNotAuthorized.Error())
		return
	}

	user, err := h.service.GetUserBySession(r.Context(), token.Value)
	if err != nil {
		slog.Error(err.Error())
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	vars := mux.Vars(r)
	resumeID := vars["resume_id"]
	slog.Info("RESUME_ID", "ID", resumeID)
	resume, err := h.service.GetUserResume(resumeID, user.ID)
	if err != nil {
		slog.Error(err.Error())
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	if err := h.service.ToggleResumeScheduling(resumeID, user.ID, resume.IsScheduled); err != nil {
		slog.Error(err.Error())
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	rr := make(map[string]string, 1)
	rr["status"] = "toggled"

	JsonResponse(w, r, http.StatusOK, rr)
}

func (h *ProfileHandler) Resumes(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		slog.Error(err.Error())
		JsonResponseErr(w, r, http.StatusForbidden, ErrNotAuthorized.Error())
		return
	}

	user, err := h.service.GetUserBySession(r.Context(), token.Value)
	if err != nil {
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	resumes, err := h.service.GetUserResumes(user.ID)
	if err != nil {
		JsonResponseErr(w, r, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	var rr ResumeResponseMany
	for _, r := range resumes {
		rr.Resumes = append(rr.Resumes, ResumeResponseSingle{
			ID:           r.ID,
			Title:        r.Title,
			AlternateUrl: r.AlternateURL,
			CreatedAt:    time.Time(r.CreatedAt),
			UpdatedAt:    time.Time(r.UpdatedAt),
			IsScheduled:  r.IsScheduled,
		})
	}

	JsonResponse(w, r, http.StatusOK, rr)
}
