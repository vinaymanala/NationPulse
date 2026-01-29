package services

import (
	"log"
	"net/http"
	"time"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
)

type DashboardService struct {
	// Add any dependencies like database connections here
	Configs *Configs
	repo    *repos.DashboardRepo
}

func NewDashboardService(configs *Configs, repo *repos.DashboardRepo) *DashboardService {
	return &DashboardService{
		Configs: configs,
		repo:    repo,
	}
}

func (ds *DashboardService) GetTopCountriesByPopulation(w http.ResponseWriter, r *http.Request) {
	log.Println("fetch top 5 populated countries")
	year := time.Now().Year()
	data, err := ds.repo.GetTopCountriesByPopulationData(year, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (ds *DashboardService) GetTopCountriesByHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("fetch top 5 health related cases in countries")
	data, err := ds.repo.GetTopCountriesByHealthData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (ds *DashboardService) GetTopCountriesByGDP(w http.ResponseWriter, r *http.Request) {
	log.Println("fetch top 5 gdp countries")
	year := time.Now().Year() - 1
	data, err := ds.repo.GetTopCountriesByGDPData(year, 10)
	if err != nil {
		log.Println("Error occured: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
