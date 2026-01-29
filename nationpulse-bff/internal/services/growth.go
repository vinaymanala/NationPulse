package services

import (
	"log"
	"net/http"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
)

type GrowthService struct {
	// Add any dependencies like database connections here
	Configs *Configs
	repo    *repos.GrowthRepo
}

func NewGrowthService(configs *Configs, repo *repos.GrowthRepo) *GrowthService {
	return &GrowthService{
		Configs: configs,
		repo:    repo,
	}
}

func (gs *GrowthService) GetGDPGrowthByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	log.Printf("fetch Gdp growth of %s\n", countryCode)
	data, err := gs.repo.GetGDPGrowthData(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (gs *GrowthService) GetPopulationGrowthByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	data, err := gs.repo.GetPopulationGrowth(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
