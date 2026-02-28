package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nationpulse-bff/internal/handlers"
	"github.com/nationpulse-bff/internal/utils"
)

type mockDashboardService struct{}

func (m *mockDashboardService) GetTopCountriesByPopulation(year, topNCountries int) (interface{}, error) {
	return []utils.TopPopulationByCountries{{CountryCode: "NG", CountryName: "Nigeria", Indicator: "Population", IndicatorCode: "POP", Year: year, Value: 200000000}}, nil
}
func (m *mockDashboardService) GetTopCountriesByHealth() (interface{}, error) {
	return []utils.TopHealthCasesByCountries{{CountryCode: "NG", CountryName: "Nigeria", Indicator: "Population", IndicatorCode: "POP", Year: 2010, Value: 200000000, SexName: "", Cause: "", UnitRange: ""}}, nil
}
func (m *mockDashboardService) GetTopCountriesByGDP(year, topNCountries int) (interface{}, error) {
	return []utils.HighestGDPCountries{}, nil
}

func TestDashboard(t *testing.T) {

	// 1. Initialize the handler with a mock service
	mux := http.NewServeMux()
	dh := handlers.NewDashboardHandler(mux, &mockDashboardService{})

	// 2. Create request and recorder
	// Test GetTopCountriesByPopulation
	var req *http.Request
	var rr *httptest.ResponseRecorder
	req = httptest.NewRequest("GET", "/api/dashboard/population", nil)
	rr = httptest.NewRecorder()

	// 3. Call handler directly
	dh.GetTopCountriesByPopulation(rr, req)
	CheckResponse(rr, t)

	// Test GetTopCountriesByHealth
	req = httptest.NewRequest("GET", "/api/dashboard/health", nil)
	rr = httptest.NewRecorder()

	dh.GetTopCountriesByHealth(rr, req)
	CheckResponse(rr, t)

	// Test GetTopCountriesByGDP
	req = httptest.NewRequest("GET", "/api/dashboard/gdp", nil)
	rr = httptest.NewRecorder()

	dh.GetTopCountriesByHealth(rr, req)
	CheckResponse(rr, t)

}

func CheckResponse(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}

	var resp utils.ApiResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.IsSuccess {
		t.Fatalf("expected success true")
	}
	dataSlice, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("unexpected data type %T", resp.Data)
	}
	if len(dataSlice) != 1 {
		t.Fatalf("expected 1 item in data, got %d", len(dataSlice))
	}
}
