package main

import (
	"context"
	"log"
	"net"
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
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//func run(ctx context.Context) {
//ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
//}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	//create the errgroup context
	g, ctx := errgroup.WithContext(ctx)

	//grpc server setup
	grpcSrv := grpc.NewServer()
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)

	// start grpc server
	g.Go(func() error {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			return err
		}
		log.Printf("gRPC health server listening at %v\n", lis.Addr())

		// goroutine to stop the gRPC server when the context is cancelled
		go func() {
			<-ctx.Done()
			grpcSrv.GracefulStop()
		}()
		return grpcSrv.Serve(lis)
	})
	// Creates a HTTP server
	srv := internals.NewServer(configs)
	httpServer := &http.Server{Addr: ":8081", Handler: srv}
	g.Go(func() error {
		log.Printf("HTTP listening on %s\n", httpServer.Addr)

		// goroutine to shutdown the HTTP server when context is cancelled
		go func() {
			<-ctx.Done()
			shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			httpServer.Shutdown(shutDownCtx)
		}()

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// Dependency Checker (Redis/Postgres)
	g.Go(func() error {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				status := grpc_health_v1.HealthCheckResponse_SERVING

				//check redis & postgres
				if err := db.Client.Ping(ctx); err != nil {
					logger.Error("DB down", zap.Error(err))
					status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
				}
				if err := rds.Client.Ping(ctx).Err(); err != nil {
					logger.Error("Redis down", zap.Error(err))
					status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
				}

				healthSrv.SetServingStatus("", status)
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Printf("System exiting with error: %v\n", err)
	}

}
