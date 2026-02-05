package services

import (
	"log"

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

func (gs *GrowthService) GetGDPGrowthByCountryCode(countryCode string) (interface{}, error) {

	log.Printf("fetch Gdp growth of %s\n", countryCode)
	return gs.repo.GetGDPGrowthData(countryCode)
}

func (gs *GrowthService) GetPopulationGrowthByCountryCode(countryCode string) (interface{}, error) {
	return gs.repo.GetPopulationGrowth(countryCode)
}
