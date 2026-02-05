package handlers

import (
	"net/http"

	. "github.com/nationpulse-bff/internal/utils"
)

type IHealthService interface {
	GetHealthByCountryCode(countryCode string) (interface{}, error)
}
type HealthHandler struct {
	mux     *http.ServeMux
	service IHealthService
}

func NewHealthHandler(mux *http.ServeMux, service IHealthService) *HealthHandler {
	return &HealthHandler{
		mux:     mux,
		service: service,
	}
}

func (hh *HealthHandler) GetHealthByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := hh.service.GetHealthByCountryCode(countryCode)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (hh *HealthHandler) RegisterRoutes() {
	hh.mux.HandleFunc("GET /country", hh.GetHealthByCountryCode)
}
