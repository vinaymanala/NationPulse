package repos

import (
	"encoding/json"
	"log"

	. "github.com/nationpulse-bff/internal/utils"
)

type EconomyRepo struct {
	Configs *Configs
}

var economyId = "economy:"

func NewEconomyRepo(configs *Configs) *EconomyRepo {
	return &EconomyRepo{
		Configs: configs,
	}
}

func (er *EconomyRepo) GetGovernmentData(countryCode string) (any, error) {
	var governmentData []EconomyData
	var cacheID = economyId + countryCode + ":government"
	data, err := GetDataFromCache(er.Configs, cacheID, &governmentData)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM geteconomygovenmentbycountry($1)`
	rows, err := er.Configs.Db.Client.Query(er.Configs.Context, sqlStatement, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var economygovernmentByCountry EconomyData

		if err := rows.Scan(
			&economygovernmentByCountry.ID,
			&economygovernmentByCountry.CountryCode,
			&economygovernmentByCountry.CountryName,
			&economygovernmentByCountry.IndicatorCode,
			&economygovernmentByCountry.Indicator,
			&economygovernmentByCountry.Year,
			&economygovernmentByCountry.Value,
			&economygovernmentByCountry.LastUpdated,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println(economygovernmentByCountry)
		governmentData = append(governmentData, economygovernmentByCountry)
	}
	if governmentData == nil {
		return governmentData, nil
	}
	marshalledData, err := json.Marshal(governmentData)
	if err != nil {
		log.Println("Error marshalling data", err)
	}
	if err := er.Configs.Cache.SetData(er.Configs.Context, cacheID, marshalledData); err != nil {
		log.Println("Error Set Cache Data", err)
	}
	return governmentData, nil
}

func (er *EconomyRepo) GetGDPData(countryCode string) (any, error) {
	var gdpData []EconomyData
	var cacheID = economyId + countryCode + "GDP"
	data, err := GetDataFromCache(er.Configs, cacheID, &gdpData)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM geteconomygdpbycountry($1)`
	rows, err := er.Configs.Db.Client.Query(er.Configs.Context, sqlStatement, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var economyGDPByCountry EconomyData

		if err := rows.Scan(
			&economyGDPByCountry.ID,
			&economyGDPByCountry.CountryCode,
			&economyGDPByCountry.CountryName,
			&economyGDPByCountry.IndicatorCode,
			&economyGDPByCountry.Indicator,
			&economyGDPByCountry.Year,
			&economyGDPByCountry.Value,
			&economyGDPByCountry.LastUpdated,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println(economyGDPByCountry)
		gdpData = append(gdpData, economyGDPByCountry)
	}
	if gdpData == nil {
		return gdpData, nil
	}
	marshalledData, err := json.Marshal(gdpData)
	if err != nil {
		log.Println("Error  marshalling data", err)
	}
	if err := er.Configs.Cache.SetData(er.Configs.Context, cacheID, marshalledData); err != nil {
		log.Println("Error Set Cache Data", err)
	}
	return gdpData, nil

}
