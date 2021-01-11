package service

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShutdownHandler struct {
	timeout time.Duration

	osSignalChan chan os.Signal
	doneChan     chan struct{}
}

type CloseFunction func(context.Context) error

// This function should be called at the beggining of the app
func NewShutdownHandler(timeout time.Duration) *ShutdownHandler {
	sh := &ShutdownHandler{
		timeout: timeout,

		osSignalChan: make(chan os.Signal, 1),
		doneChan:     make(chan struct{}),
	}

	signal.Notify(sh.osSignalChan, syscall.SIGTERM, syscall.SIGINT)

	return sh
}

func (sh ShutdownHandler) WaitShutdown(closer CloseFunction) {
	go func() {
		<-sh.osSignalChan
		log.Println("Received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), sh.timeout)
		defer cancel()

		if err := closer(ctx); err != nil {
			log.Println("Failed to shutdown gracefully")
		} else {
			log.Println("Shutdown gracefully")
		}

		sh.doneChan <- struct{}{}
	}()

	<-sh.doneChan
}

func (sh ShutdownHandler) Close() {
	close(sh.osSignalChan)
	close(sh.doneChan)
}
