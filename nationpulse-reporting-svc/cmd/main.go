package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/service"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/store"
	. "github.com/vinaymanala/nationpulse-reporting-svc/internal/types"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Load environment variables from .env for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load; relying on environment variables")
	}

	cfg := config.Load()

	rds := store.NewRedis(cfg)
	db := store.NewPgClient(ctx, cfg)

	configs := &Configs{
		Cfg:   cfg,
		Ctx:   ctx,
		DB:    db,
		Cache: rds,
	}
	exportService := service.NewExportService(configs)
	exportService.Serve()

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

	// Dependency Checker (Postgres)
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
				//check redis & postgres
				if err := db.Client.Ping(ctx); err != nil {
					log.Printf("DB down: %v\n", err)
					status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
				}
				if err := rds.Client.Ping(ctx).Err(); err != nil {
					log.Printf("Redis down: %v\n", err)
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
