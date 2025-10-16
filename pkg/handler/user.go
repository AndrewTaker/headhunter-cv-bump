package handler

import (
	"fmt"
	"log"
	"net/http"
	"pkg/auth"
	"pkg/service"
)

type UserHandler struct {
	userService service.UserService
	auth        *auth.AuthRepository
}

func NewUserHandler(us service.UserService, auth *auth.AuthRepository) *UserHandler {
	return &UserHandler{userService: us, auth: auth}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("sess")
	if err != nil {
		http.Error(w, "Could not retrieve cookie", http.StatusBadRequest)
		log.Println("cookie err", err)
		return
	}

	userID := h.auth.GetUserByToken(token.Value)
	// user, err := h.userService.GetUser("60645454")
	// if err != nil {
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<strong>I validated your request as user %s</strong>", userID)))
}
