package handlers

import (
	"net/http"

	"github.com/nationpulse-bff/internal/services"
)

type UtilsHandler struct {
	mux     *http.ServeMux
	service *services.UtilsService
}

func NewUtilsHandler(mux *http.ServeMux, service *services.UtilsService) *UtilsHandler {
	return &UtilsHandler{
		mux:     mux,
		service: service,
	}
}

func (uh *UtilsHandler) RegisterRoutes() {
	uh.mux.HandleFunc("POST /permissions", uh.service.GetUserPermissions)
	uh.mux.HandleFunc("POST /reports/publish", uh.service.PublishExportRequest)
	uh.mux.HandleFunc("GET /reports/subscribe/event", uh.service.SubscribeExportResponse)

}
