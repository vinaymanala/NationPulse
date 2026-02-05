package handlers

import (
	"net/http"

	. "github.com/nationpulse-bff/internal/utils"
)

type IEconomyService interface {
	GetEconomyGovernmentDataByCountryCode(countryCode string) (interface{}, error)
	GetEconomyGDPByCountryCode(countryCode string) (interface{}, error)
}
type EconomyHandler struct {
	mux     *http.ServeMux
	service IEconomyService
}

func NewEconomyHandler(mux *http.ServeMux, service IEconomyService) *EconomyHandler {
	return &EconomyHandler{
		mux:     mux,
		service: service,
	}
}

func (eh *EconomyHandler) GetEconomyGovernmentDataByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := eh.service.GetEconomyGovernmentDataByCountryCode(countryCode)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (eh *EconomyHandler) GetEconomyGDPByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := eh.service.GetEconomyGDPByCountryCode(countryCode)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (eh *EconomyHandler) RegisterRoutes() {
	eh.mux.HandleFunc("GET /governmentdata/country", eh.GetEconomyGovernmentDataByCountryCode)
	eh.mux.HandleFunc("GET /gdp/country", eh.GetEconomyGDPByCountryCode)

}
