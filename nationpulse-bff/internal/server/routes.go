package internals

import (
	"fmt"
	"net/http"

	"github.com/nationpulse-bff/internal/handlers"
	"github.com/nationpulse-bff/internal/repos"
	"github.com/nationpulse-bff/internal/services"
	"github.com/nationpulse-bff/internal/utils"
)

func handleDashboardRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling Dashboard Route..")
}
func handleHealthRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling Health Route..")
}
func handlePopulationRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling Population Route..")
}
func handleEconomyRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling Economy Route..")
}
func handleGrowthRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling PerformanceGrowth Route..")
}

func addRoutes(configs *utils.Configs, muxes *ServerMux) {

	adminService := services.NewAdminService(configs, repos.NewAdminRepo(configs))
	userService := services.NewUserService(configs, repos.NewUserRepo(configs))
	utilsService := services.NewUtilsService(configs, repos.NewUtilsRepo(configs))
	dashboardService := services.NewDashboardService(configs, repos.NewDashboardRepo(configs))
	populationService := services.NewPopulationService(configs, repos.NewPopulationRepo(configs))
	healthService := services.NewHealthService(configs, repos.NewHealthRepo(configs))
	economyService := services.NewEconomyService(configs, repos.NewEconomyRepo(configs))
	growthService := services.NewGrowthService(configs, repos.NewGrowthRepo(configs))

	ah := handlers.NewAdminHandler(muxes.AdminMux, adminService)
	ah.RegisterRoutes()

	uh := handlers.NewUserHandler(muxes.UserMux, userService)
	uh.RegisterRoutes()

	utilsS := handlers.NewUtilsHandler(muxes.UtilsMux, utilsService)
	utilsS.RegisterRoutes()

	muxes.DashboardMux.HandleFunc("GET /", handleDashboardRoute)
	dh := handlers.NewDashboardHandler(muxes.DashboardMux, dashboardService)
	dh.RegisterRoutes()

	muxes.PopulationMux.HandleFunc("GET /", handlePopulationRoute)
	ph := handlers.NewPopulationHandler(muxes.PopulationMux, populationService)
	ph.RegisterRoutes()

	muxes.HealthMux.HandleFunc("GET /", handleHealthRoute)
	hh := handlers.NewHealthHandler(muxes.HealthMux, healthService)
	hh.RegisterRoutes()

	muxes.EconomyMux.HandleFunc("GET /", handleEconomyRoute)
	eh := handlers.NewEconomyHandler(muxes.EconomyMux, economyService)
	eh.RegisterRoutes()

	muxes.GrowthMux.HandleFunc("GET /", handleGrowthRoute)
	gh := handlers.NewGrowthHandler(muxes.GrowthMux, growthService)
	gh.RegisterRoutes()
}
