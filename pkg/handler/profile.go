package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkg/auth"
	"pkg/headhunter"
	"pkg/service"
	"strconv"
)

type ProfileHandler struct {
	service *service.SqliteService
	auth    *auth.AuthRepository
}

func NewProfileHandler(s *service.SqliteService, a *auth.AuthRepository) *ProfileHandler {
	return &ProfileHandler{
		service: s,
		auth:    a,
	}
}

func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("cookie err", err)
		return
	}

	userID := h.auth.GetUserByToken(token.Value)
	user, err := h.service.GetUser(userID)
	if err != nil {
		return
	}

	resumes, err := h.service.GetUserResumes(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("db err", err)
		return
	}

	log.Println(user, resumes)
}

func (h *ProfileHandler) GetResumes(w http.ResponseWriter, r *http.Request) {
	var resumes []headhunter.Resume

	for i := range 10 {
		s := strconv.Itoa(i)
		resumes = append(resumes, headhunter.Resume{
			ID:           s,
			AlternateURL: fmt.Sprintf("https://localhost.com/api/%s", s),
			Title:        fmt.Sprintf("title for resumes %s", s),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(resumes)
	w.Write(data)

	// token, err := r.Cookie("sess")
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	log.Println("cookie err", err)
	// 	return
	// }
	//
	// userID := h.auth.GetUserByToken(token.Value)
	//
	// dbToken, err := h.tokenService.GetTokenByUserID(r.Context(), userID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	log.Println("db err", err)
	// 	return
	// }
	//
	// oauth2Token := dbToken.ToOauth2Token()
	//
	// client := headhunter.HHOauthConfig.Client(r.Context(), oauth2Token)
	// resp, err := client.Get("https://api.hh.ru/resumes/mine")
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	log.Println("db err", err)
	// 	return
	// }
	// defer resp.Body.Close()
	//
	// resumes, err := h.resumeService.GetUserResumes(userID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	log.Println("db err", err)
	// 	return
	// }
	//
	// log.Println(resumes)
}
