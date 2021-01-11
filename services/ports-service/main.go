package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis"
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

	repo, repoCloser, err := createRepository(cfg)
	if err != nil {
		log.Printf("Failed to create repository. Error: %s", err)
		os.Exit(1)
	}
	portService := service.NewPort(repo)

	grpcServer, err := StartGrpcServerAsync(cfg, portService)
	if err != nil {
		log.Printf("Failed to start gRPC server. Error: %s", err)
		os.Exit(1)
	}

	// closeFunc will be called by shutdown handler when it receives a shutdown signal
	closeFunc := func(ctx context.Context) error {
		log.Printf("Closing service resources...")
		grpcServer.Stop()
		if repoCloser != nil {
			repoCloser.Close()
		}
		// add here other resources which needs to be close
		log.Printf("Closed all resources!")

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

func createRepository(cfg *config.Config) (repository.Repository, io.Closer, error) {
	if cfg.RedisEndpoint == "" {
		return repository.NewMemoryStorage(), nil, nil
	}

	log.Printf("Using redis repository on %s", cfg.RedisEndpoint)

	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisEndpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rc.Ping().Result()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to redis on %s. Error:%s", cfg.RedisEndpoint, err)
	}

	return repository.NewRedis(rc), rc, nil
}
