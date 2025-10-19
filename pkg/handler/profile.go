package handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"pkg/auth"
	"pkg/headhunter"
	"pkg/service"

	"github.com/starfederation/datastar-go/datastar"
)

type Store struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

type ProfileHandler struct {
	userService  service.UserService
	tokenService service.TokenService
	auth         *auth.AuthRepository
	tmpl         *template.Template
}

func NewProfileHandler(
	ts service.TokenService,
	us service.UserService,
	auth *auth.AuthRepository,
	tmpl *template.Template,
) *ProfileHandler {
	return &ProfileHandler{tokenService: ts, userService: us, auth: auth, tmpl: tmpl}
}

func (h *ProfileHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": "not logged in"})
		log.Println("cookie err", err)
		return
	}

	userID := h.auth.GetUserByToken(token.Value)
	user, err := h.userService.GetUser(userID)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": err.Error()})
		return
	}

	h.tmpl.ExecuteTemplate(w, "index", map[string]any{"User": user})
}

func (h *ProfileHandler) GetResumes(w http.ResponseWriter, r *http.Request) {
	store := &Store{}
	if err := datastar.ReadSignals(r, store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("sess")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": "not logged in"})
		log.Println("cookie err", err)
		return
	}

	userID := h.auth.GetUserByToken(token.Value)
	// user, err := h.userService.GetUser(userID)
	// if err != nil {
	// 	h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": err.Error()})
	// 	return
	// }

	dbToken, err := h.tokenService.GetTokenByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": "not logged in"})
		log.Println("db err", err)
		return
	}

	oauth2Token := dbToken.ToOauth2Token()

	client := headhunter.HHOauthConfig.Client(r.Context(), oauth2Token)
	resp, err := client.Get("https://api.hh.ru/resumes/mine")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": "not logged in"})
		log.Println("db err", err)
		return
	}
	defer resp.Body.Close()

	type hhResumesResponse struct {
		Items []headhunter.Resume `json:"items"`
	}
	var hhr hhResumesResponse
	err = json.NewDecoder(resp.Body).Decode(&hhr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "index", map[string]string{"Error": "not logged in"})
		log.Println("db err", err)
		return
	}

	sse := datastar.NewSSE(w, r)

	var t bytes.Buffer
	err = h.tmpl.ExecuteTemplate(&t, "resumes", map[string]any{"Resumes": hhr.Items})
	if err != nil {
		log.Println("template err", err)
	}
	sse.PatchElements(t.String())
}
