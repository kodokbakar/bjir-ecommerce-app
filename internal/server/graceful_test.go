package server

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

type fakeShutdowner struct {
	called atomic.Bool
	err    error
}

func (f *fakeShutdowner) Shutdown(ctx context.Context) error {
	f.called.Store(true)
	return f.err
}

func TestGracefulShutdownCallsShutdownAndClosers(t *testing.T) {
	shutdowner := &fakeShutdowner{}

	var shutdownStarted bool
	var firstClosed bool
	var secondClosed bool

	err := GracefulShutdown(context.Background(), shutdowner, GracefulShutdownOptions{
		Timeout: 100 * time.Millisecond,
		OnShutdownStart: func() {
			shutdownStarted = true
		},
		CloseFuncs: []CloseFunc{
			func() error {
				firstClosed = true
				return nil
			},
			func() error {
				secondClosed = true
				return nil
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !shutdowner.called.Load() {
		t.Fatal("expected shutdown to be called")
	}

	if !shutdownStarted {
		t.Fatal("expected shutdown start hook to be called")
	}

	if !firstClosed {
		t.Fatal("expected first close func to be called")
	}

	if !secondClosed {
		t.Fatal("expected second close func to be called")
	}
}

func TestGracefulShutdownIgnoresHTTPServerClosed(t *testing.T) {
	shutdowner := &fakeShutdowner{
		err: http.ErrServerClosed,
	}

	err := GracefulShutdown(context.Background(), shutdowner, GracefulShutdownOptions{
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("expected http.ErrServerClosed to be ignored, got %v", err)
	}
}
