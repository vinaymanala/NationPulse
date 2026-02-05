package services

import (
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

func (ps *PopulationService) GetPopulationByCountryCode(countryCode string) (interface{}, error) {
	return ps.repo.GetPopulationByCountryData(countryCode)
}
