package utils

import (
	"context"
	"time"

	"github.com/nationpulse-bff/internal/config"
	"github.com/nationpulse-bff/internal/kafka"
	"github.com/nationpulse-bff/internal/store"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Configs struct {
	Db                  *store.PgClient
	Cache               *store.Redis
	Context             context.Context
	Kafka               *kafka.Kafka
	Logger              *zap.Logger
	Metrics             *prometheus.Registry
	MetricHttpRequests  *prometheus.CounterVec
	MetricHttpDurations *prometheus.HistogramVec
	Cfg                 config.Config
}

type Filter struct {
	Headers []string   `json:"headers"`
	Records [][]string `json:"records"`
	Query   string     `json:"query"`
}

type ExportApiMessageRequest struct {
	ExportID           string `json:"exportID"`
	UserID             int    `json:"userID"`
	Filters            Filter `json:"filters"`
	RequestTableString string `json:"requestTableString"`
	RequestCountryCode string `json:"requestCountryCode"`
}

type ExportApiMessageResponse struct {
	Status     string `json:"Status"`
	StatusCode int    `json:"StatusCode"`
	Data       any    `json:"Data"`
}

type ApiResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data      any  `json:"data"`
	Error     any  `json:"error"`
}
type TopPopulationByCountries struct {
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	Indicator     string  `json:"indicator"`
	IndicatorCode string  `json:"indicator_code"`
	Year          int     `json:"year"`
	Value         float64 `json:"value"`
}

type TopHealthCasesByCountries struct {
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	Indicator     string  `json:"indicator"`
	IndicatorCode string  `json:"indicator_code"`
	Year          int     `json:"year"`
	Value         float64 `json:"value"`
	SexName       string  `json:"sex_name"`
	Cause         string  `json:"cause"`
	UnitRange     string  `json:"unitRange"`
}

type PopulationData struct {
	ID            int       `json:"id"`
	CountryCode   string    `json:"country_code"`
	CountryName   string    `json:"country_name"`
	IndicatorCode string    `json:"indicator_code"`
	Indicator     string    `json:"indicator"`
	SexCode       string    `json:"sex_code"`
	SexName       string    `json:"sex_name"`
	Age           string    `json:"age"`
	Year          int       `json:"year"`
	Value         float64   `json:"value"`
	LastUpdated   time.Time `json:"last_updated"`
}

type HealthData struct {
	ID            int       `json:"id"`
	CountryCode   string    `json:"country_code"`
	CountryName   string    `json:"country_name"`
	IndicatorCode string    `json:"indicator_code"`
	Indicator     string    `json:"indicator"`
	SexCode       string    `json:"sex_code"`
	SexName       string    `json:"sex_name"`
	Cause         string    `json:"cause"`
	UnitRange     string    `json:"unitRange"`
	Year          int       `json:"year"`
	Value         float64   `json:"value"`
	LastUpdated   time.Time `json:"last_updated"`
}

type HighestGDPCountries struct {
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	Indicator     string  `json:"indicator"`
	IndicatorCode string  `json:"indicator_code"`
	Year          string  `json:"year"`
	Value         float64 `json:"value"`
}

type EconomyData struct {
	ID            int       `json:"id"`
	CountryCode   string    `json:"country_code"`
	CountryName   string    `json:"country_name"`
	IndicatorCode string    `json:"indicator_code"`
	Indicator     string    `json:"indicator"`
	Year          string    `json:"year"`
	Value         float64   `json:"value"`
	LastUpdated   time.Time `json:"last_updated"`
}

type GrowthData struct {
	ID            int       `json:"id"`
	CountryCode   string    `json:"country_code"`
	CountryName   string    `json:"country_name"`
	IndicatorCode string    `json:"indicator_code"`
	Indicator     string    `json:"indicator"`
	Year          string    `json:"year"`
	Value         float64   `json:"value"`
	LastUpdated   time.Time `json:"last_updated"`
}

type UserPermissions struct {
	Name            string `json:"username"`
	Email           string `json:"email"`
	RoleId          int    `json:"role_id"`
	RoleName        string `json:"role_name"`
	RoleDescription string `json:"role_description"`
	ModuleID        int    `json:"module_id"`
	ModuleName      string `json:"module_name"`
	ModuleValue     int    `json:"module_value"`
	PermissionID    int    `json:"permission_id"`
	PermissionName  string `json:"permission_name"`
	PermissionValue int    `json:"permission_value"`
}

type UpdatePermissions struct {
	UserID      int   `json:"user_id"`
	RoleID      int   `json:"role_id"`
	Modules     []int `json:"modules"`
	Permissions []int `json:"permissions"`
}

type Users struct {
	ID        int       `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpatedAt  time.Time `json:"updated_at"`
}
