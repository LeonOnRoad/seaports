package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	pkgService "company.com/seaports/pkg/service"
	pbPort "company.com/seaports/proto/src/api/port"

	"company.com/seaports/services/ports-api/config"
	"company.com/seaports/services/ports-api/server"
	"company.com/seaports/services/ports-api/service"
)

func main() {
	cfg := config.LoadConfig()

	sh := pkgService.NewShutdownHandler(cfg.ShutdownTimeout)
	defer sh.Close()

	serverResources, err := createServerResources(cfg)
	if err != nil {
		log.Printf("Failed to create server resources. Error: %s", err)
		os.Exit(1)
	}

	httpServer := server.StartAsync(cfg.Port, serverResources)

	closeFunc := func(ctx context.Context) error {
		// this function will be called by shutdown handler when it receives a shutdown signal
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Println("Failed to close http server gracefully")
			return err
		}

		// add here other resources which needs to be close

		return nil
	}

	sh.WaitShutdown(closeFunc) // blocking
}

func createServerResources(cfg *config.Config) (*server.Resources, error) {
	portsConn, err := grpc.Dial(cfg.PortsServiceEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	portsClient := pbPort.NewPortsServiceClient(portsConn)
	portService := service.NewPort(portsClient)

	return &server.Resources{
		PortService: portService,
	}, nil
}
