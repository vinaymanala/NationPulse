package repos

import (
	"encoding/json"
	"log"

	. "github.com/nationpulse-bff/internal/utils"
)

var populationID = "population:"

type PopulationRepo struct {
	Configs *Configs
}

func NewPopulationRepo(configs *Configs) *PopulationRepo {
	return &PopulationRepo{
		Configs: configs,
	}
}

func (pr *PopulationRepo) GetPopulationByCountryData(countryCode string) (any, error) {
	var populationByCountries []PopulationData
	var cacheID = populationID + countryCode
	data, err := GetDataFromCache(pr.Configs, cacheID, &populationByCountries)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM getpopulationofcountry($1)`
	rows, err := pr.Configs.Db.Client.Query(pr.Configs.Context, sqlStatement, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var populationByCountry PopulationData
		if err := rows.Scan(
			&populationByCountry.ID,
			&populationByCountry.CountryCode,
			&populationByCountry.CountryName,
			&populationByCountry.IndicatorCode,
			&populationByCountry.Indicator,
			&populationByCountry.SexCode,
			&populationByCountry.SexName,
			&populationByCountry.Age,
			&populationByCountry.Year,
			&populationByCountry.Value,
			&populationByCountry.LastUpdated,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println(populationByCountry)
		populationByCountries = append(populationByCountries, populationByCountry)
	}
	if populationByCountries == nil {
		return populationByCountries, nil
	}
	marshalledData, err := json.Marshal(populationByCountries)
	if err != nil {
		log.Println("Error marshalling data", marshalledData)
	}
	if err := pr.Configs.Cache.SetData(pr.Configs.Context, cacheID, marshalledData); err != nil {
		log.Println("Error Set Cache Data", err)
	}
	return populationByCountries, nil
}
