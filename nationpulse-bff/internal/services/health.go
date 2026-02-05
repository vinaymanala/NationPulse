package services

import (
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

func (hs *HealthService) GetHealthByCountryCode(countryCode string) (interface{}, error) {
	// log.Printf("fetch health of %s\n", countryCode)
	return hs.repo.GetHealthData(countryCode)
}
