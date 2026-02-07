package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/joho/godotenv"
	"github.com/vinaymanala/nationpulse-data-ingestion-svc/internal/config"
	"github.com/vinaymanala/nationpulse-data-ingestion-svc/internal/service"
	"github.com/vinaymanala/nationpulse-data-ingestion-svc/internal/store"
	. "github.com/vinaymanala/nationpulse-data-ingestion-svc/internal/types"
	"github.com/vinaymanala/nationpulse-data-ingestion-svc/pb"
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

	// r := gin.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pg := *store.NewPgClient(ctx, cfg)

	Configs := &Configs{
		Cfg: cfg,
		Ctx: ctx,
		DB:  &pg,
	}

	dataIngestionSvc := service.NewDataIngestionSvc(Configs)
	// dataIngestionSvc.Serve()
	// os.Exit(1)

	//create the errgroup context
	g, ctx := errgroup.WithContext(ctx)

	//grpc server setup
	grpcSrv := grpc.NewServer()
	// protobuf : Data ingestion svc
	pb.RegisterDataIngestionServer(grpcSrv, dataIngestionSvc)

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
				if err := pg.Client.Ping(ctx); err != nil {
					log.Printf("DB down: %v\n", err)
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
