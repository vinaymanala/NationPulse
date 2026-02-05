package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	. "github.com/nationpulse-bff/internal/utils"
)

type IAdminService interface {
	GetUserPermissions(userID string) (interface{}, error)
	SetUserPermissions(data UpdatePermissions) error
	GetUsers() (interface{}, error)
}

type AdminHandler struct {
	mux     *http.ServeMux
	service IAdminService
}

func NewAdminHandler(mux *http.ServeMux, service IAdminService) *AdminHandler {
	return &AdminHandler{
		mux:     mux,
		service: service,
	}
}

func (ah *AdminHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req struct{ UserID string }
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	data, err := ah.service.GetUserPermissions(req.UserID)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (ah *AdminHandler) SetUserPermissions(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data UpdatePermissions
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = ah.service.SetUserPermissions(data)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (ah *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	data, err := ah.service.GetUsers()
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, false, err.Error())
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (ah *AdminHandler) RegisterRoutes() {
	ah.mux.HandleFunc("POST /getUserPermissions", ah.GetUserPermissions)
	ah.mux.HandleFunc("POST /setUserPermissions", ah.SetUserPermissions)
	ah.mux.HandleFunc("GET /getUsers", ah.GetUsers)
}
