package main_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nationpulse-bff/internal/config"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load; relying on environment variables")
	}

	// Allow skipping DB setup for fast unit tests: set SKIP_DB=1
	if os.Getenv("SKIP_DB") == "1" {
		code := m.Run()
		os.Exit(code)
	}

	// Load configuration from environment
	cfg := config.Load()

	// initialize pgxpool
	ctx := context.Background()
	connStr := "postgres://" + cfg.PostgresUser + ":" + cfg.PostgresPass + "@" + cfg.PostgresAddr + "/" + cfg.PostgresName + "?sslmode=disable"
	fmt.Println("Test DB:", connStr)

	var err error
	testPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// check if db is actually reachable
	if err := testPool.Ping(ctx); err != nil {
		log.Fatalf("Database unreachable: %v\n", err)
	}

	//2: Run tests
	code := m.Run()

	//3: Teardown
	testPool.Close()

	os.Exit(code)
}
