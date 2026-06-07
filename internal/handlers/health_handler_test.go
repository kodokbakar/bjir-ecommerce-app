package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck_ReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetShuttingDown(false)

	router := gin.New()
	router.GET("/health", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"success":true`) {
		t.Fatalf("expected success true, got: %s", w.Body.String())
	}
}

func TestHealthCheck_ReturnsServiceUnavailableDuringShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetShuttingDown(true)
	defer SetShuttingDown(false)

	router := gin.New()
	router.GET("/health", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"code":"service_unavailable"`) {
		t.Fatalf("expected service_unavailable code, got: %s", w.Body.String())
	}
}
