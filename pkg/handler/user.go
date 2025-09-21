package handler

import (
	"net/http"
	"pkg/service"
)

type UserHandler struct {
	userService service.UserService
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Response) {
	user, err := h.userService.GetUser("60645454")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(user.FirstName))
}
