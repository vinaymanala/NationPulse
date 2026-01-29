package repos

import (
	"encoding/json"
	"log"
	"strconv"

	. "github.com/nationpulse-bff/internal/utils"
)

type UtilsRepo struct {
	Configs *Configs
}

func NewUtilsRepo(configs *Configs) *UtilsRepo {
	return &UtilsRepo{
		Configs: configs,
	}
}

var utilsID = "utils:"

func (ur *UtilsRepo) GetPermissions(userID string) (any, error) {
	var userPermissions []UserPermissions
	var permissions []int
	data, err := GetDataFromCache(ur.Configs, utilsID+"permissions:"+userID, &permissions)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return *data, nil
	}

	sqlStatement := `SELECT * FROM get_user_permissions($1);`
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Fatal("Error converting userId to int")
		return nil, err
	}

	permissionsData, err := FetchPermissionsFromDB(ur.Configs, sqlStatement, userPermissions, id)
	if err != nil {
		log.Println("Error fetching permissions from DB", err)
	}

	// if err := ur.Configs.Cache.SetData(ur.Configs.Context, utilsID+"modulePermissions:"+userID, permissionsData); err != nil {
	// 	log.Println("Error Set Cache Data for route level permissions", err)
	// }

	for _, permission := range permissionsData {
		permissions = append(permissions, permission.ModuleValue)
	}
	if permissions == nil {
		return permissions, nil
	}
	marshalledData, err := json.Marshal(permissions)
	if err != nil {
		log.Println("Error marshalling data", err)
	}
	if err := ur.Configs.Cache.SetData(ur.Configs.Context, utilsID+"permissions:"+userID, marshalledData); err != nil {
		log.Println("Error Set Cache Data for user permissions", err)
	}
	return permissions, nil
}
