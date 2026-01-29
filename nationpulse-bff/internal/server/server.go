package internals

import (
	"net/http"

	"github.com/nationpulse-bff/internal/middlewares"
	"github.com/nationpulse-bff/internal/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServerMux struct {
	UserMux       *http.ServeMux
	AdminMux      *http.ServeMux
	DashboardMux  *http.ServeMux
	PopulationMux *http.ServeMux
	HealthMux     *http.ServeMux
	EconomyMux    *http.ServeMux
	GrowthMux     *http.ServeMux
	UtilsMux      *http.ServeMux
}

func (sm *ServerMux) NewServerMuxes() *ServerMux {
	return &ServerMux{
		UserMux:       http.NewServeMux(),
		AdminMux:      http.NewServeMux(),
		DashboardMux:  http.NewServeMux(),
		PopulationMux: http.NewServeMux(),
		HealthMux:     http.NewServeMux(),
		EconomyMux:    http.NewServeMux(),
		GrowthMux:     http.NewServeMux(),
		UtilsMux:      http.NewServeMux(),
	}
}

func groupRoutePrefix(prefix string, mux *http.ServeMux) http.Handler {
	return http.StripPrefix(prefix, mux)
}

func NewServer(configs *utils.Configs) http.Handler {

	rootMux := http.NewServeMux()
	newMux := &ServerMux{}
	muxes := newMux.NewServerMuxes()

	rootMux.Handle("/metrics", promhttp.Handler())
	rootMux.Handle("/api/u/",
		middlewares.DefaultMiddlewares(configs, groupRoutePrefix("/api/u", muxes.UserMux)))
	rootMux.Handle("/api/uu/",
		middlewares.WithAuthMiddlewares(configs, groupRoutePrefix("/api/uu", muxes.UtilsMux)))
	rootMux.Handle("/api/a/",
		middlewares.WithAuthMiddlewares(configs, groupRoutePrefix("/api/a", muxes.AdminMux)))
	rootMux.Handle("/api/dashboard/",
		middlewares.DefaultMiddlewares(configs, groupRoutePrefix("/api/dashboard", muxes.DashboardMux)))
	rootMux.Handle("/api/health/",
		middlewares.WithAuthMiddlewares(configs, groupRoutePrefix("/api/health", muxes.HealthMux)))
	rootMux.Handle("/api/population/",
		middlewares.WithAuthMiddlewares(configs, groupRoutePrefix("/api/population", muxes.PopulationMux)))
	rootMux.Handle("/api/economy/",
		middlewares.WithAuthMiddlewares(configs, groupRoutePrefix("/api/economy", muxes.EconomyMux)))
	rootMux.Handle("/api/growth/",
		middlewares.WithAuthMiddlewares(configs, groupRoutePrefix("/api/growth", muxes.GrowthMux)))

	addRoutes(configs, muxes)

	return rootMux
}
