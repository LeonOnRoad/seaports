package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pkgService "company.com/seaports/pkg/service"
	pbPort "company.com/seaports/proto/src/api/port"

	"company.com/seaports/services/ports-service/config"
	"company.com/seaports/services/ports-service/repository"
	"company.com/seaports/services/ports-service/service"
)

func main() {
	cfg := config.LoadConfig()

	sh := pkgService.NewShutdownHandler(cfg.ShutdownTimeout)
	defer sh.Close()

	repo := repository.NewMemoryStorage()
	portService := service.NewPort(repo)

	grpcServer, err := StartGrpcServerAsync(cfg, portService)
	if err != nil {
		//panic(err)
		log.Printf("Failed to start gRPC server. Error: %s", err)
		os.Exit(1)
	}

	closeFunc := func(ctx context.Context) error {
		grpcServer.Stop()

		// add here other resources which need to be close

		return nil
	}

	sh.WaitShutdown(closeFunc) // blocking
}

func StartGrpcServerAsync(cfg *config.Config, portService *service.Port) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	pbPort.RegisterPortsServiceServer(s, portService)

	// Start grpc service
	go func() {
		log.Println("Start PortService gRPC server")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
		log.Println("Stoped PortService gRPC server")
	}()
	return s, nil
}
