package main

import (
	"context"
	"io"
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

	serverResources, serverClosers, err := createServerResources(cfg)
	if err != nil {
		log.Printf("Failed to create server resources. Error: %s", err)
		os.Exit(1)
	}

	httpServer := server.StartAsync(cfg.Port, serverResources)

	// closeFunc will be called by shutdown handler when it receives a shutdown signal
	closeFunc := func(ctx context.Context) error {
		log.Printf("Closing service resources...")
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Println("Failed to close http server gracefully")
			return err
		}

		for _, closer := range serverClosers {
			if closer != nil {
				closer.Close()
			}
		}

		// add here other resources which needs to be close

		log.Printf("Closed all resources!")

		return nil
	}

	sh.WaitShutdown(closeFunc) // blocking
}

func createServerResources(cfg *config.Config) (*server.Resources, []io.Closer, error) {
	closerRes := make([]io.Closer, 0) // add here all needed resources that needs to be closed when app exists
	portsConn, err := grpc.Dial(cfg.PortsServiceEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	portsClient := pbPort.NewPortsServiceClient(portsConn)
	portService := service.NewPort(portsClient)

	closerRes = append(closerRes, portsConn) // is nice that grpc connection to be closed when app exists

	return &server.Resources{
		PortService: portService,
	}, closerRes, nil
}
