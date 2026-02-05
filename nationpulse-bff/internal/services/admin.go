package services

import (
	"errors"

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

func (as *AdminService) GetUserPermissions(userID string) (interface{}, error) {
	if userID == "" {
		return nil, errors.New("Invalid UserID")
	}
	return as.repo.GetUserPermissions(userID)
}

func (as *AdminService) SetUserPermissions(data UpdatePermissions) error {
	// upate the database with new permissions
	return as.repo.SetUserPermissions(data)
}

func (as *AdminService) GetUsers() (interface{}, error) {
	return as.repo.GetUsers()
}
