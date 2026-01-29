package services

import (
	"log"
	"net/http"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
)

type PopulationService struct {
	// Add any dependencies like database connections here
	Configs *Configs
	repo    *repos.PopulationRepo
}

func NewPopulationService(configs *Configs, repo *repos.PopulationRepo) *PopulationService {
	return &PopulationService{
		Configs: configs,
		repo:    repo,
	}
}

func (ps *PopulationService) GetPopulationByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	log.Printf("fetch population of %s\n", countryCode)
	data, err := ps.repo.GetPopulationByCountryData(countryCode)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
