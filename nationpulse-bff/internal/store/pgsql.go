package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nationpulse-bff/internal/config"
)

type PgClient struct {
	Client *pgxpool.Pool
}

var (
	pgInstance *PgClient
	pgOnce     sync.Once
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewPgClient(ctx context.Context, cfg config.Config) *PgClient {
	// pgHost := cfg.PostgresHost
	pgName := cfg.PostgresName
	pgPass := cfg.PostgresPass
	pgUser := cfg.PostgresUser
	pgAddr := cfg.PostgresAddr
	// connStr := "postgres://postgres:postgres@localhost:5432/nationPulseDB?sslmode=disable"
	connStr := "postgres://" + pgUser + ":" + pgPass + "@" + pgAddr + "/" + pgName + "?sslmode=disable"
	fmt.Println(connStr)
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connStr)
		if err != nil {
			fmt.Printf("Error occured while connecting database: %s\n", err)
			panic(err)
		}
		pgInstance = &PgClient{Client: db}
	})
	fmt.Println("Connected to Postgres database successfully")
	return pgInstance
}

func (pg *PgClient) Ping(ctx context.Context) error {
	return pg.Client.Ping(ctx)
}

func (pg *PgClient) GetUser(ctx context.Context, user *User) (*User, error) {
	sqlStatement := `SELECT * FROM get_user($1, $2);`
	row := pg.Client.QueryRow(ctx, sqlStatement, user.Name, user.Email)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	fmt.Printf("Result: id: %s, user:%s, email:%s \n ", user.ID, user.Name, user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}
