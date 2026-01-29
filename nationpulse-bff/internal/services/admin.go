package services

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
)

type AdminService struct {
	// Add any dependencies like database connections here
	Configs *Configs
	repo    *repos.AdminRepo
}

func NewAdminService(configs *Configs, repo *repos.AdminRepo) *AdminService {
	return &AdminService{
		Configs: configs,
		repo:    repo,
	}
}

func (as *AdminService) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
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

	data, err := as.repo.GetUserPermissions(req.UserID)
	if err != nil {
		http.Error(w, "Error fetching permissions: "+err.Error(), http.StatusBadRequest)
		WriteJSON(w, http.StatusBadRequest, nil, false, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (as *AdminService) SetUserPermissions(w http.ResponseWriter, r *http.Request) {
	// upate the database with new permissions
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
	err = as.repo.SetUserPermissions(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (as *AdminService) GetUsers(w http.ResponseWriter, r *http.Request) {

	data, err := as.repo.GetUsers()
	if err != nil {
		http.Error(w, "Error fetching all users"+err.Error(), http.StatusBadRequest)
		WriteJSON(w, http.StatusBadRequest, nil, false, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, data, true, nil)
}
