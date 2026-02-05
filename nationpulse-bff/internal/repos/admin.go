package repos

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
	. "github.com/nationpulse-bff/internal/utils"
)

type AdminRepo struct {
	Configs *Configs
}

func NewAdminRepo(configs *Configs) *AdminRepo {
	return &AdminRepo{
		Configs: configs,
	}
}

func (ar *AdminRepo) GetUserPermissions(userID string) (interface{}, error) {
	var userPermissions []UserPermissions
	var permissions []int

	data, err := GetDataFromCache(ar.Configs, "utils:permissions:"+userID, &userPermissions)
	if err != nil {
		log.Println("Cache Get Failed. Trying DB.")
	} else {
		return data, nil
	}

	sqlStatement := `SELECT * FROM get_user_permissions($1);`
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Println("Error converting userId to int")
		return nil, err
	}

	permissionsData, err := FetchPermissionsFromDB(ar.Configs, sqlStatement, userPermissions, id)
	if err != nil {
		log.Println("Error fetching permissions from DB", err)
	}

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
	if err := ar.Configs.Cache.SetData(ar.Configs.Context, "utils:permissions:"+userID, marshalledData); err != nil {
		log.Println("Error Set Cache Data for user permissions", err)
	}
	return permissions, nil
}

func (ar *AdminRepo) SetUserPermissions(updatePermissions UpdatePermissions) error {
	// upate the database with new permissions
	tx, err := ar.Configs.Db.Client.Begin(ar.Configs.Context)
	if err != nil {
		return err
	}
	// rollback if not committed
	defer tx.Rollback(ar.Configs.Context)
	// Call stored procedure
	_, err = tx.Exec(ar.Configs.Context,
		"Select  update_user_permissions($1, $2, $3, $4)",
		updatePermissions.UserID,
		updatePermissions.RoleID,
		updatePermissions.Modules,
		updatePermissions.Permissions)
	if err != nil {
		log.Println("failed to update permissions:")
		return err
	}

	if err := tx.Commit(ar.Configs.Context); err != nil {
		return err
	}

	err = ar.Configs.Cache.DelData(ar.Configs.Context, "utils:permissions:"+strconv.Itoa(updatePermissions.UserID))
	if err != nil {
		log.Printf("Failed to invalidate cache for user %d: %v", updatePermissions.UserID, err)
	}

	return nil
}

func (ar *AdminRepo) GetUsers() (interface{}, error) {
	var users []Users
	sqlStatement := `SELECT * from users`
	rows, err := ar.Configs.Db.Client.Query(ar.Configs.Context, sqlStatement)
	if err != nil {
		log.Println("Error fetching users from DB", err)
		return nil, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByPos[Users])
	if err != nil {
		log.Printf("Error scanning rows: %v\n", err)
	}

	return users, nil
}
