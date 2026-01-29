package handlers

import (
	"net/http"

	"github.com/nationpulse-bff/internal/services"
)

type HealthHandler struct {
	mux     *http.ServeMux
	service *services.HealthService
}

func NewHealthHandler(mux *http.ServeMux, service *services.HealthService) *HealthHandler {
	return &HealthHandler{
		mux:     mux,
		service: service,
	}
}

func (hh *HealthHandler) RegisterRoutes() {
	hh.mux.HandleFunc("GET /country", hh.service.GetHealthByCountryCode)
}
