package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

type CloseFunc func() error

type GracefulShutdownOptions struct {
	Timeout         time.Duration
	OnShutdownStart func()
	CloseFuncs      []CloseFunc
}

func GracefulShutdown(parentCtx context.Context, shutdowner Shutdowner, options GracefulShutdownOptions) error {
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	if options.Timeout <= 0 {
		options.Timeout = 30 * time.Second
	}

	if options.OnShutdownStart != nil {
		options.OnShutdownStart()
	}

	shutdownCtx, cancel := context.WithTimeout(parentCtx, options.Timeout)
	defer cancel()

	err := shutdowner.Shutdown(shutdownCtx)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	if err != nil {
		log.Printf("failed to shutdown http server gracefully: %v", err)
	}

	for _, closeFunc := range options.CloseFuncs {
		if closeFunc == nil {
			continue
		}

		if closeErr := closeFunc(); closeErr != nil {
			log.Printf("failed to close resource during shutdown: %v", closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}

	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		log.Printf("graceful shutdown timeout reached after %s", options.Timeout)
		if err == nil {
			err = shutdownCtx.Err()
		}
	}

	return err
}
