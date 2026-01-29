package handlers

import (
	"net/http"

	"github.com/nationpulse-bff/internal/services"
)

type AdminHandler struct {
	mux     *http.ServeMux
	service *services.AdminService
}

func NewAdminHandler(mux *http.ServeMux, service *services.AdminService) *AdminHandler {
	return &AdminHandler{
		mux:     mux,
		service: service,
	}
}

func (ah *AdminHandler) RegisterRoutes() {
	ah.mux.HandleFunc("POST /getUserPermissions", ah.service.GetUserPermissions)
	ah.mux.HandleFunc("POST /setUserPermissions", ah.service.SetUserPermissions)
	ah.mux.HandleFunc("GET /getUsers", ah.service.GetUsers)

}
