package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nationpulse-bff/internal/config"
	"github.com/nationpulse-bff/internal/kafka"
	internals "github.com/nationpulse-bff/internal/server"
	"github.com/nationpulse-bff/internal/store"
	"github.com/nationpulse-bff/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

//func run(ctx context.Context) {
//ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
//}

func main() {
	ctx := context.Background()
	// Load environment variables from .env for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load; relying on environment variables")
	}

	// Load configuration from environment
	cfg := config.Load()

	// Initialize redis store
	rds := store.NewRedis(cfg, 15*time.Minute)
	// Initialize postgres client
	db := store.NewPgClient(ctx, cfg)
	// Initialize kafka svc
	k := kafka.NewKafka(ctx, cfg)
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	//Initialize metrics
	metrics := prometheus.NewRegistry()
	httpRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration in seconds",
		},
		[]string{"method", "path"},
	)

	prometheus.MustRegister(httpRequests, httpDuration)

	configs := &utils.Configs{
		Db:                  db,
		Cache:               rds,
		Context:             ctx,
		Kafka:               k,
		Logger:              logger,
		Metrics:             metrics,
		MetricHttpRequests:  httpRequests,
		MetricHttpDurations: httpDuration,
		Cfg:                 cfg,
	}

	defer rds.Client.Close()
	defer db.Client.Close()

	// Creates a HTTP server
	srv := internals.NewServer(configs)
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: srv,
	}
	fmt.Printf("Starting up..")
	go func() {
		log.Printf("listening to %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("error listening and serving: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	httpServer.Shutdown(ctx)
}
