package handler

import (
	"net/http"
	"pkg/headhunter"
	"pkg/service"
)

type UserHandler struct {
	userService service.UserService
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Response) {
	ctx := r.Context()

	client := headhunter.NewHHClient(ctx)
}
