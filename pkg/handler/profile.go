package handler

import (
	"log/slog"
	"net/http"
	"pkg/service"
	"time"
)

type ProfileHandler struct {
	service *service.SqliteService
}

func NewProfileHandler(s *service.SqliteService) *ProfileHandler {
	return &ProfileHandler{service: s}
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
