package handlers

import (
	"net/http"

	. "github.com/nationpulse-bff/internal/utils"
)

type IGrowthHandler interface {
	GetGDPGrowthByCountryCode(countryCode string) (interface{}, error)
	GetPopulationGrowthByCountryCode(countryCode string) (interface{}, error)
}
type GrowthHandler struct {
	mux     *http.ServeMux
	service IGrowthHandler
}

func NewGrowthHandler(mux *http.ServeMux, service IGrowthHandler) *GrowthHandler {
	return &GrowthHandler{
		mux:     mux,
		service: service,
	}
}
func (gh *GrowthHandler) GetGDPGrowthByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := gh.service.GetGDPGrowthByCountryCode(countryCode)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (gh *GrowthHandler) GetPopulationGrowthByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := gh.service.GetPopulationGrowthByCountryCode(countryCode)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
func (gh *GrowthHandler) RegisterRoutes() {
	gh.mux.HandleFunc("GET /gdp/country", gh.GetGDPGrowthByCountryCode)
	gh.mux.HandleFunc("GET /population/country", gh.GetPopulationGrowthByCountryCode)
}
