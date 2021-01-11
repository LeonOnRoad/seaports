package main

import (
	"context"

	pkgService "company.com/seaports/pkg/service"

	"company.com/seaports/services/ports-service/config"
)

func main() {
	cfg := config.LoadConfig()
	sh := pkgService.NewShutdownHandler(cfg.ShutdownTimeout)
	defer sh.Close()

	closeFunc := func(ctx context.Context) error {

		// add here other resources which need to be close

		return nil
	}

	sh.WaitShutdown(closeFunc) // blocking
}
