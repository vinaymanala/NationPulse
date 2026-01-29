package repos

import (
	"encoding/json"
	"log"

	. "github.com/nationpulse-bff/internal/utils"
)

type GrowthRepo struct {
	Configs *Configs
}

func NewGrowthRepo(configs *Configs) *GrowthRepo {
	return &GrowthRepo{
		Configs: configs,
	}
}

var growthID = "growth:"

func (gr *GrowthRepo) GetGDPGrowthData(countryCode string) (any, error) {
	var gdpGrowthData []GrowthData
	var cacheID = growthID + countryCode + ":GDP"
	data, err := GetDataFromCache(gr.Configs, cacheID, &gdpGrowthData)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM getgrowthgdpofcountry($1)`
	rows, err := gr.Configs.Db.Client.Query(gr.Configs.Context, sqlStatement, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var growthgdpDataByCountry GrowthData

		if err := rows.Scan(
			&growthgdpDataByCountry.ID,
			&growthgdpDataByCountry.CountryCode,
			&growthgdpDataByCountry.CountryName,
			&growthgdpDataByCountry.IndicatorCode,
			&growthgdpDataByCountry.Indicator,
			&growthgdpDataByCountry.Year,
			&growthgdpDataByCountry.Value,
			&growthgdpDataByCountry.LastUpdated,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println(growthgdpDataByCountry)
		gdpGrowthData = append(gdpGrowthData, growthgdpDataByCountry)
	}

	if gdpGrowthData == nil {
		return gdpGrowthData, nil
	}
	marshalledData, err := json.Marshal(gdpGrowthData)
	if err != nil {
		log.Println("Error marshalling data", err)
	}
	gr.Configs.Cache.SetData(gr.Configs.Context, cacheID, marshalledData)
	return gdpGrowthData, nil
}

func (gr *GrowthRepo) GetPopulationGrowth(countryCode string) (any, error) {
	var populationGrowthData []GrowthData
	var cacheID = growthID + countryCode + ":population"
	data, err := GetDataFromCache(gr.Configs, cacheID, &populationGrowthData)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM getgrowthpopnbycountry($1)`
	rows, err := gr.Configs.Db.Client.Query(gr.Configs.Context, sqlStatement, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var populationGrowthByCountry GrowthData

		if err := rows.Scan(
			&populationGrowthByCountry.ID,
			&populationGrowthByCountry.CountryCode,
			&populationGrowthByCountry.CountryName,
			&populationGrowthByCountry.IndicatorCode,
			&populationGrowthByCountry.Indicator,
			&populationGrowthByCountry.Year,
			&populationGrowthByCountry.Value,
			&populationGrowthByCountry.LastUpdated,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println(populationGrowthByCountry)
		populationGrowthData = append(populationGrowthData, populationGrowthByCountry)
	}
	if populationGrowthData == nil {
		return populationGrowthData, nil
	}
	marshalledData, err := json.Marshal(populationGrowthData)
	if err != nil {
		log.Println("Error marshalling data", err)
	}
	if err := gr.Configs.Cache.SetData(gr.Configs.Context, cacheID, marshalledData); err != nil {
		log.Println("Error Set Cache Data", err)
	}
	return populationGrowthData, nil

}
