package services

import (
	"log"

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

func (ds *DashboardService) GetTopCountriesByPopulation(year, topNCountries int) (interface{}, error) {
	log.Println("fetch top 5 populated countries")
	return ds.repo.GetTopCountriesByPopulationData(year, topNCountries)
}

func (ds *DashboardService) GetTopCountriesByHealth() (interface{}, error) {
	log.Println("fetch top 5 health related cases in countries")
	return ds.repo.GetTopCountriesByHealthData()

}

func (ds *DashboardService) GetTopCountriesByGDP(year, topNCountries int) (interface{}, error) {
	log.Println("fetch top 5 gdp countries")
	return ds.repo.GetTopCountriesByGDPData(year, topNCountries)
}
