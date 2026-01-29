package handlers

import (
	"net/http"

	"github.com/nationpulse-bff/internal/services"
)

type GrowthHandler struct {
	mux     *http.ServeMux
	service *services.GrowthService
}

func NewGrowthHandler(mux *http.ServeMux, service *services.GrowthService) *GrowthHandler {
	return &GrowthHandler{
		mux:     mux,
		service: service,
	}
}

func (gh *GrowthHandler) RegisterRoutes() {
	gh.mux.HandleFunc("GET /gdp/country", gh.service.GetGDPGrowthByCountryCode)
	gh.mux.HandleFunc("GET /population/country", gh.service.GetPopulationGrowthByCountryCode)
}
