package services

import (
	"net/http"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
)

type EconomyService struct {
	// Add any dependencies like database connections here
	Configs *Configs
	repo    *repos.EconomyRepo
}

func NewEconomyService(configs *Configs, repo *repos.EconomyRepo) *EconomyService {
	return &EconomyService{
		Configs: configs,
		repo:    repo,
	}
}

func (es *EconomyService) GetEconomyGovernmentDataByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	// log.Printf("Economy GovermentData of %s\n", countryCode)
	data, err := es.repo.GetGovernmentData(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (es *EconomyService) GetEconomyGDPByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	// log.Printf("Economy GDP of %s\n", countryCode)
	data, err := es.repo.GetGDPData(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
