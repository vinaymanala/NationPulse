package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func WriteJSON(w http.ResponseWriter, status int, data any, success bool, err any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var response = &ApiResponse{
		Data:      data,
		IsSuccess: success,
		Error:     err,
	}
	// log.Printf("Message: %s: %s, isSuccess: %t, Error:%v \n", success, err)
	return json.NewEncoder(w).Encode(response)
}

func GetUserDetailsFromCache(r *http.Request, configs *Configs) {
	fmt.Printf("Form Details: %s\n", r.Form.Get("userID"))
}

func GetDataFromCache[T any](configs *Configs, key string, mappedStruct T) (*T, error) {
	//var zero T
	data, err := configs.Cache.GetData(configs.Context, key)
	// fmt.Println("CACHE DATA:", data, err)
	if err != nil {
		log.Printf("Error fetching data from cache %s\n", err)
		return nil, errors.New("error fetching data from cache")
	}
	if err := json.Unmarshal([]byte(data), &mappedStruct); err != nil {
		log.Println("Error unmarshalling data from cache.")
		return nil, errors.New("error unmarshalling data from cache")
	}
	// fmt.Println("==================================")
	// fmt.Println("Unmarshal data", &mappedStruct)
	fmt.Println("Fetched Data from Cache!!")
	return &mappedStruct, nil
}

func checkModulePermission(permissions []UserPermissions, moduleID int) bool {
	// log.Printf("HEALTH_ID %d", HEALTH_ID)
	log.Printf("Permissions LOG: %v", permissions)
	for _, p := range permissions {
		log.Printf("p.ModuleID: %d, moduleID: %d", p.ModuleValue, moduleID)
		if p.ModuleValue == moduleID {
			return true
		}
	}
	return false
}

func HasPermissions(requestPath string, permissions *[]UserPermissions) bool {
	log.Printf("PERMISSIONS ARG %v", *permissions)
	switch {
	case strings.HasPrefix(requestPath, ADMIN_PERMISSION):
		return checkModulePermission(*permissions, PERMISSION_ID)
	case strings.HasPrefix(requestPath, DASHBOARD) || strings.HasPrefix(requestPath, PERMISSION):
		return checkModulePermission(*permissions, DASHBOARD_ID)
	case strings.HasPrefix(requestPath, POPULATION):
		return checkModulePermission(*permissions, POPULATION_ID)
	case strings.HasPrefix(requestPath, HEALTH):
		log.Println("Touched Health Case")
		return checkModulePermission(*permissions, HEALTH_ID)
	case strings.HasPrefix(requestPath, ECONOMY):
		return checkModulePermission(*permissions, ECONOMY_ID)
	case strings.HasPrefix(requestPath, GROWTH):
		return checkModulePermission(*permissions, GROWTH_ID)
	case strings.HasPrefix(requestPath, REPORTING):
		return checkModulePermission(*permissions, REPORTING_ID)
	default:
		return false
	}
}

func FetchPermissionsFromDB(configs *Configs, sqlStatement string, userPermissions []UserPermissions, id int) ([]UserPermissions, error) {

	rows, err := configs.Db.Client.Query(configs.Context, sqlStatement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// for rows.Next() {
	// 	var userPermission UserPermissions
	// if err := rows.Scan(
	// 	&userPermission.Name,
	// 	&userPermission.Email,
	// 	&userPermission.RoleId,
	// 	&userPermission.RoleName,
	// 	&userPermission.RoleDescription,
	// 	&userPermission.ModuleID,
	// 	&userPermission.ModuleName,
	// 	&userPermission.ModuleValue,
	// 	&userPermission.PermissionID,
	// 	&userPermission.PermissionName,
	// 	&userPermission.PermissionValue,
	// );
	data, err := pgx.CollectRows(rows, pgx.RowToStructByPos[UserPermissions])
	if err != nil {
		log.Fatalf("Error scanning a row: %v\n", err)
		return nil, err
	}
	// userPermissions = append(userPermissions, userPermission)
	return data, nil

}

func GetModulePermissionsFromCache(configs *Configs, userID int, key string, permissions []UserPermissions, w http.ResponseWriter, r *http.Request) ([]UserPermissions, error) {
	cacheData, err := configs.Cache.GetData(configs.Context, key)
	if err != nil {
		// cache miss, fetch db
		sqlStatement := `SELECT * FROM get_user_permissions($1);`
		data, err := FetchPermissionsFromDB(configs, sqlStatement, permissions, userID)
		log.Println("Fetching permissions from DB.")
		if err != nil {
			log.Println("Error: UnAuthorized")
			WriteJSON(w, http.StatusUnauthorized, nil, false, "Cannot fetch permissions, User not authroized.")
		}
		marshalledData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data")
			return nil, err
		}
		if err := configs.Cache.SetData(configs.Context, "utils:modulePermissions:"+strconv.Itoa(userID), marshalledData); err != nil {
			log.Println("Error Set Cache Data for route level permissions", err)
			return data, nil
		}
		return data, nil
	}
	if err := json.Unmarshal([]byte(cacheData), &permissions); err != nil {
		log.Println("Error unmarshalling permissions dataa", err)
		return nil, err
	}
	return permissions, nil
}

func getJSONTags(s any) []string {
	t := reflect.TypeOf(s)
	fields := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		fields[i] = jsonTag
	}

	return fields
}

func GetQueryAndHeaders(req *ExportApiMessageRequest) error {

	table := req.RequestTableString
	var query string
	var headers []string
	switch table {
	case "population":
		query = `SELECT * FROM getpopulationofcountry($1)`
		headers = getJSONTags(PopulationData{})
	case "health":
		query = `SELECT * FROM gethealthcasesbycountry($1)`
		headers = getJSONTags(HealthData{})
	case "economy:gov":
		query = `SELECT * FROM geteconomygovenmentbycountry($1)`
		headers = getJSONTags(EconomyData{})
	case "economy:gdp":
		query = `SELECT * FROM geteconomygdpbycountry($1)`
		headers = getJSONTags(EconomyData{})
	case "growth:gdp":
		query = `SELECT * FROM getgrowthgdpofcountry($1)`
		headers = getJSONTags(GrowthData{})
	case "growth:population":
		query = `SELECT * FROM getgrowthpopnbycountry($1)`
		headers = getJSONTags(GrowthData{})
	}

	req.Filters.Query = query
	req.Filters.Headers = headers

	return nil
}

func ReadMessages(r *kafka.Reader, ctx context.Context,
	callbackFn func(message kafka.Message) error) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Kafka BFF Reader shutting down...")
			return
		default:
		}
		readerCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()
		message, err := r.ReadMessage(readerCtx)
		if err != nil {
			log.Println("Error Reading message: ", err)
			continue
		}
		// fmt.Println("MESSAGE========>", string(message.Value))
		if err := callbackFn(message); err != nil {
			log.Println("Error occured processing the message", err)
			continue
		}
		if err := r.CommitMessages(ctx, message); err != nil {
			log.Println("Error committing", err)
		}

	}
}
