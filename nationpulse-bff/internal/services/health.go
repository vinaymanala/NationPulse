package services

import (
	"net/http"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
)

type HealthService struct {
	// Add any dependencies like database connections here
	Configs *Configs
	repo    *repos.HealthRepo
}

func NewHealthService(configs *Configs, repo *repos.HealthRepo) *HealthService {
	return &HealthService{
		Configs: configs,
		repo:    repo,
	}
}

func (hs *HealthService) GetHealthByCountryCode(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("countryCode")
	// log.Printf("fetch health of %s\n", countryCode)
	data, err := hs.repo.GetHealthData(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}
