package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupCORSRouter(config CORSConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSWithConfig(config))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	return router
}

func TestCORSHeadersPresent(t *testing.T) {
	router := setupCORSRouter(CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:       12 * time.Hour,
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://frontend.example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected allow origin *, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if !strings.Contains(w.Header().Get("Access-Control-Allow-Methods"), "GET") {
		t.Fatalf("expected allow methods to contain GET, got %s", w.Header().Get("Access-Control-Allow-Methods"))
	}

	if !strings.Contains(w.Header().Get("Access-Control-Allow-Headers"), "Authorization") {
		t.Fatalf("expected allow headers to contain Authorization, got %s", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if !strings.Contains(w.Header().Get("Access-Control-Allow-Headers"), "Content-Type") {
		t.Fatalf("expected allow headers to contain Content-Type, got %s", w.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestCORSPreflightRequest(t *testing.T) {
	router := setupCORSRouter(CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:       12 * time.Hour,
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://frontend.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected allow origin *, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if !strings.Contains(w.Header().Get("Access-Control-Allow-Methods"), "POST") {
		t.Fatalf("expected allow methods to contain POST, got %s", w.Header().Get("Access-Control-Allow-Methods"))
	}

	if !strings.Contains(w.Header().Get("Access-Control-Allow-Headers"), "Authorization") {
		t.Fatalf("expected allow headers to contain Authorization, got %s", w.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestCORSRestrictedOriginAllowed(t *testing.T) {
	router := setupCORSRouter(CORSConfig{
		AllowOrigins: []string{"https://frontend.example.com"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:       12 * time.Hour,
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://frontend.example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "https://frontend.example.com" {
		t.Fatalf("expected restricted origin to be allowed, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Vary") != "Origin" {
		t.Fatalf("expected Vary Origin header, got %s", w.Header().Get("Vary"))
	}
}

func TestCORSRestrictedOriginBlocked(t *testing.T) {
	router := setupCORSRouter(CORSConfig{
		AllowOrigins: []string{"https://frontend.example.com"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:       12 * time.Hour,
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://evil.example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("expected blocked origin to have no allow-origin header, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestLoadCORSConfigFromEnv(t *testing.T) {
	t.Setenv("CORS_ALLOW_ORIGINS", "https://a.example.com,https://b.example.com")
	t.Setenv("CORS_ALLOW_METHODS", "GET,POST,OPTIONS")
	t.Setenv("CORS_ALLOW_HEADERS", "Authorization,Content-Type,X-Request-ID")
	t.Setenv("CORS_MAX_AGE", "24h")

	config := LoadCORSConfigFromEnv()

	if len(config.AllowOrigins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(config.AllowOrigins))
	}

	if config.AllowOrigins[0] != "https://a.example.com" {
		t.Fatalf("expected first origin https://a.example.com, got %s", config.AllowOrigins[0])
	}

	if len(config.AllowMethods) != 3 {
		t.Fatalf("expected 3 methods, got %d", len(config.AllowMethods))
	}

	if len(config.AllowHeaders) != 3 {
		t.Fatalf("expected 3 headers, got %d", len(config.AllowHeaders))
	}

	if config.MaxAge != 24*time.Hour {
		t.Fatalf("expected max age 24h, got %v", config.MaxAge)
	}
}
