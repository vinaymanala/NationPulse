package handlers

import (
	"net/http"

	"github.com/nationpulse-bff/internal/services"
)

type PopulationHandler struct {
	mux     *http.ServeMux
	service *services.PopulationService
}

func NewPopulationHandler(mux *http.ServeMux, service *services.PopulationService) *PopulationHandler {
	return &PopulationHandler{
		mux:     mux,
		service: service,
	}
}

func (ph *PopulationHandler) RegisterRoutes() {
	ph.mux.HandleFunc("GET /country", ph.service.GetPopulationByCountryCode)
}
