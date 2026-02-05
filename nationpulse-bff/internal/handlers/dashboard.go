package handlers

import (
	"net/http"
	"time"

	. "github.com/nationpulse-bff/internal/utils"
)

type IDashboardService interface {
	GetTopCountriesByPopulation(year, topNCountries int) (interface{}, error)
	GetTopCountriesByHealth() (interface{}, error)
	GetTopCountriesByGDP(year, topNCountries int) (interface{}, error)
}
type DashboardHandler struct {
	// Add any dependencies like services here
	mux     *http.ServeMux
	service IDashboardService
}

func NewDashboardHandler(mux *http.ServeMux, service IDashboardService) *DashboardHandler {
	return &DashboardHandler{
		mux:     mux,
		service: service,
	}
}

func (dh *DashboardHandler) GetTopCountriesByPopulation(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Year()
	data, err := dh.service.GetTopCountriesByPopulation(year, 10)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (dh *DashboardHandler) GetTopCountriesByHealth(w http.ResponseWriter, r *http.Request) {
	data, err := dh.service.GetTopCountriesByHealth()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (dh *DashboardHandler) GetTopCountriesByGDP(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Year() - 1
	data, err := dh.service.GetTopCountriesByGDP(year, 10)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (dh *DashboardHandler) RegisterRoutes() {
	// Register dashboard-related routes here
	//dh.mux.HandleFunc("GET /", handleDashboardRoute)
	dh.mux.HandleFunc("GET /population", dh.GetTopCountriesByPopulation)
	dh.mux.HandleFunc("GET /health", dh.GetTopCountriesByHealth)
	dh.mux.HandleFunc("GET /gdp", dh.GetTopCountriesByGDP)

}
