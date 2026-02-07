package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vinaymanala/nationpulse-cron-data-ingestion-schedular-svc/internal/config"
	s "github.com/vinaymanala/nationpulse-cron-data-ingestion-schedular-svc/internal/scheduler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	sch := s.New(cfg)
	sch.Start()
	defer sch.Stop()

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

	if err := g.Wait(); err != nil {
		log.Printf("System exiting with error: %v\n", err)
	}
}
