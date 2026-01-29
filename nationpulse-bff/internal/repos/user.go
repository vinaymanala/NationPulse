package repos

import (
	"log"

	"github.com/nationpulse-bff/internal/store"
	. "github.com/nationpulse-bff/internal/utils"
)

type UserRepo struct {
	Configs *Configs
}

func NewUserRepo(configs *Configs) *UserRepo {
	return &UserRepo{
		Configs: configs,
	}
}

func (ur *UserRepo) GetUserDetails(user *store.User) (*store.User, error) {
	sqlStatement := `SELECT * FROM get_user($1, $2);`
	rowOne := ur.Configs.Db.Client.QueryRow(ur.Configs.Context, sqlStatement, user.Name, user.Email)

	err := rowOne.Scan(&user.ID, &user.Name, &user.Email)
	// fmt.Printf("Result: id: %s, user:%s, email:%s \n ", user.ID, user.Name, user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) GetPermissions(user *store.User) (any, error) {
	var userPermissions []UserPermissions
	sqlStatement := `SELECT * FROM get_user_permissions($1);`
	rows, err := ur.Configs.Db.Client.Query(ur.Configs.Context, sqlStatement, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userPermission UserPermissions
		if err := rows.Scan(
			&userPermission.Name,
			&userPermission.Email,
			&userPermission.RoleId,
			&userPermission.RoleName,
			&userPermission.RoleDescription,
			&userPermission.ModuleID,
			&userPermission.ModuleName,
			&userPermission.ModuleValue,
			&userPermission.PermissionID,
			&userPermission.PermissionName,
			&userPermission.PermissionValue,
		); err != nil {
			log.Fatalf("Error scanning a row: %v\n", err)
			return nil, err
		}
		userPermissions = append(userPermissions, userPermission)
	}

	var permissions []int
	for _, permission := range userPermissions {
		permissions = append(permissions, permission.PermissionValue)
	}
	// fmt.Println("PERMISSIONS", permissions)
	return permissions, nil
}
