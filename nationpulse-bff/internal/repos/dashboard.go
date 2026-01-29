package repos

import (
	"encoding/json"
	"log"

	. "github.com/nationpulse-bff/internal/utils"
)

type DashboardRepo struct {
	Configs *Configs
}

var dashboardID = "dashboard:"

func NewDashboardRepo(configs *Configs) *DashboardRepo {
	return &DashboardRepo{
		Configs: configs,
	}
}

func (dr *DashboardRepo) GetTopCountriesByPopulationData(currentYear int, top_countries int) (any, error) {
	var topCountriesByPopulation []TopPopulationByCountries
	data, err := GetDataFromCache(dr.Configs, dashboardID+"population", &topCountriesByPopulation)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	// Cache nil, checking DB
	sqlStatement := `SELECT * FROM get_perfgrowthpopulation_dashboard($1, $2);`
	rows, err := dr.Configs.Db.Client.Query(dr.Configs.Context, sqlStatement, currentYear, top_countries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var topPopulatedCountry TopPopulationByCountries
		if err := rows.Scan(
			&topPopulatedCountry.CountryCode,
			&topPopulatedCountry.CountryName,
			&topPopulatedCountry.Indicator,
			&topPopulatedCountry.IndicatorCode,
			&topPopulatedCountry.Year,
			&topPopulatedCountry.Value); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		//fmt.Println(topPopulatedCountry)
		topCountriesByPopulation = append(topCountriesByPopulation, topPopulatedCountry)
	}
	if topCountriesByPopulation == nil {
		return topCountriesByPopulation, nil
	}
	marshalledData, err := json.Marshal(topCountriesByPopulation)
	if err != nil {
		log.Println("Error marshaling data to cache")
	}
	if err := dr.Configs.Cache.SetData(dr.Configs.Context, dashboardID+"population", marshalledData); err != nil {
		log.Println("Error set cache data:", err)
	}
	// fmt.Println("Call Successfull")
	return topCountriesByPopulation, nil
}

func (dr *DashboardRepo) GetTopCountriesByHealthData() (any, error) {
	var topCountriesByHealth []TopHealthCasesByCountries
	if data, err := GetDataFromCache(dr.Configs, dashboardID+"health", &topCountriesByHealth); err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}
	sqlStatement := `SELECT * FROM get_healthstatus_dashboard()`
	rows, err := dr.Configs.Db.Client.Query(dr.Configs.Context, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var topHealthDatByCountry TopHealthCasesByCountries
		if err := rows.Scan(
			&topHealthDatByCountry.CountryCode,
			&topHealthDatByCountry.CountryName,
			&topHealthDatByCountry.Indicator,
			&topHealthDatByCountry.IndicatorCode,
			&topHealthDatByCountry.Year,
			&topHealthDatByCountry.Value,
			&topHealthDatByCountry.SexName,
			&topHealthDatByCountry.Cause,
			&topHealthDatByCountry.UnitRange,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		//fmt.Println(topHealthDatByCountry)
		topCountriesByHealth = append(topCountriesByHealth, topHealthDatByCountry)
	}
	if topCountriesByHealth == nil {
		return topCountriesByHealth, nil
	}
	marshalledData, err := json.Marshal(topCountriesByHealth)
	if err != nil {
		log.Println("Error marshalling data", err)
	}
	if err := dr.Configs.Cache.SetData(dr.Configs.Context, dashboardID+"healthdata", marshalledData); err != nil {
		log.Println("Error Set Cache data", err)
	}
	return topCountriesByHealth, nil
}

func (dr *DashboardRepo) GetTopCountriesByGDPData(currentYear int, topNCountries int) (any, error) {
	var highestGDPCountries []HighestGDPCountries
	if data, err := GetDataFromCache(dr.Configs, dashboardID+"GDPdata", &highestGDPCountries); err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}

	sqlStatement := `SELECT * FROM get_highest_gdppercapita_countries_dashboard($1, $2)`
	rows, err := dr.Configs.Db.Client.Query(dr.Configs.Context, sqlStatement, currentYear, topNCountries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var highestGDPCountry HighestGDPCountries
		if err := rows.Scan(
			&highestGDPCountry.CountryCode,
			&highestGDPCountry.CountryName,
			&highestGDPCountry.Indicator,
			&highestGDPCountry.IndicatorCode,
			&highestGDPCountry.Year,
			&highestGDPCountry.Value); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		// fmt.Println("COMing Dat", highestGDPCountry)
		highestGDPCountries = append(highestGDPCountries, highestGDPCountry)
	}

	if highestGDPCountries == nil {
		// fmt.Println("highestGDPCountries is null=============")
		return highestGDPCountries, nil
	}
	marshalledData, err := json.Marshal(highestGDPCountries)
	if err != nil {
		log.Println("Error marshalling data", err)
	}
	if err := dr.Configs.Cache.SetData(dr.Configs.Context, dashboardID+"GDPdata", marshalledData); err != nil {
		log.Println("Error Set Cache Data", err)
	}
	return highestGDPCountries, nil
}
