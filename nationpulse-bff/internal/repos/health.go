package repos

import (
	"encoding/json"
	"log"

	. "github.com/nationpulse-bff/internal/utils"
)

type HealthRepo struct {
	Configs *Configs
}

func NewHealthRepo(configs *Configs) *HealthRepo {
	return &HealthRepo{
		Configs: configs,
	}
}

var healthID = "health:"

func (hr *HealthRepo) GetHealthData(countryCode string) (any, error) {
	var healthData []HealthData
	var cacheID = healthID + countryCode
	data, err := GetDataFromCache(hr.Configs, cacheID, &healthData)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM gethealthcasesbycountry($1)`
	rows, err := hr.Configs.Db.Client.Query(hr.Configs.Context, sqlStatement, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var healthDataByCountry HealthData

		if err := rows.Scan(
			&healthDataByCountry.ID,
			&healthDataByCountry.CountryCode,
			&healthDataByCountry.CountryName,
			&healthDataByCountry.IndicatorCode,
			&healthDataByCountry.Indicator,
			&healthDataByCountry.SexCode,
			&healthDataByCountry.SexName,
			&healthDataByCountry.Cause,
			&healthDataByCountry.UnitRange,
			&healthDataByCountry.Year,
			&healthDataByCountry.Value,
			&healthDataByCountry.LastUpdated,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println(healthDataByCountry)
		healthData = append(healthData, healthDataByCountry)
	}
	if healthData == nil {
		return healthData, nil
	}
	marshalledData, err := json.Marshal(healthData)
	if err != nil {
		log.Println("Error marshalling data", marshalledData)
	}
	if err := hr.Configs.Cache.SetData(hr.Configs.Context, cacheID, marshalledData); err != nil {
		log.Println("Error Set Cache Data", err)
	}
	return healthData, nil
}
