package handlers

import (
	"net/http"

	"github.com/nationpulse-bff/internal/services"
)

type UserHandler struct {
	mux     *http.ServeMux
	service *services.UserService
}

func NewUserHandler(mux *http.ServeMux, service *services.UserService) *UserHandler {
	return &UserHandler{
		mux:     mux,
		service: service,
	}
}

func (uh *UserHandler) RegisterRoutes() {
	uh.mux.HandleFunc("POST /signin", uh.service.HandleLogin)
	uh.mux.HandleFunc("POST /signout", uh.service.HandleLogout)
	uh.mux.HandleFunc("POST /token/refresh", uh.service.HandleRefreshToken)

}
