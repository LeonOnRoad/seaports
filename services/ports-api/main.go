package main

import (
	"context"
	"log"
	"time"

	pkgService "company.com/seaports/pkg/service"
	"company.com/seaports/services/ports-api/server"
)

func main() {
	sh := pkgService.NewShutdownHandler(10 * time.Second)
	defer sh.Close()

	httpServer := server.StartAsync(8080)

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

	sh.WaitShutdown(closeFunc)
}
