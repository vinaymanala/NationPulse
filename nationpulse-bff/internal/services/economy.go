package services

import (
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

func (es *EconomyService) GetEconomyGovernmentDataByCountryCode(countryCode string) (interface{}, error) {
	// log.Printf("Economy GovermentData of %s\n", countryCode)
	return es.repo.GetGovernmentData(countryCode)

}

func (es *EconomyService) GetEconomyGDPByCountryCode(countryCode string) (interface{}, error) {
	// log.Printf("Economy GDP of %s\n", countryCode)
	return es.repo.GetGDPData(countryCode)
}
