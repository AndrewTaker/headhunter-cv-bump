package handler

import (
	"log"
	"net/http"
	"pkg/service"
)

type ProfileHandler struct {
	service *service.SqliteService
}

func NewProfileHandler(s *service.SqliteService) *ProfileHandler {
	return &ProfileHandler{service: s}
}

func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		JsonResponseErr(w, r, http.StatusForbidden, ErrNotAuthorized.Error())
		return
	}

	user, err := h.service.GetUserBySession(r.Context(), token.Value)
	if err != nil {
		log.Println(err)
		return
	}

	resumes, err := h.service.GetUserResumes(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("db err", err)
		return
	}

	log.Println(user, resumes)
}
