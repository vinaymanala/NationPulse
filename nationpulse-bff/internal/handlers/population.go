package handlers

import (
	"net/http"

	. "github.com/nationpulse-bff/internal/utils"
)

type IPopulationService interface {
	GetPopulationByCountryCode(countryCode string) (interface{}, error)
}

type PopulationHandler struct {
	mux     *http.ServeMux
	service IPopulationService
}

func NewPopulationHandler(mux *http.ServeMux, service IPopulationService) *PopulationHandler {
	return &PopulationHandler{
		mux:     mux,
		service: service,
	}
}

func (ph *PopulationHandler) GetPopulationByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := ph.service.GetPopulationByCountryCode(countryCode)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
func (ph *PopulationHandler) RegisterRoutes() {
	ph.mux.HandleFunc("GET /country", ph.GetPopulationByCountryCode)
}
