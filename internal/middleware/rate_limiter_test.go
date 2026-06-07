package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRateLimitRouter(config RateLimitConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RateLimitWithLimiter(NewRateLimiter(config)))

	router.GET("/api/v1/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.GET("/api/v1/cart", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	return router
}

func testRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Public: RateLimitRule{
			Limit:  2,
			Window: time.Minute,
		},
		Auth: RateLimitRule{
			Limit:  3,
			Window: time.Minute,
		},
		Login: RateLimitRule{
			Limit:  1,
			Window: time.Minute,
		},
	}
}

func TestRateLimitHeadersPresent(t *testing.T) {
	router := setupRateLimitRouter(testRateLimitConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("X-RateLimit-Limit") != "2" {
		t.Fatalf("expected limit header 2, got %s", w.Header().Get("X-RateLimit-Limit"))
	}

	if w.Header().Get("X-RateLimit-Remaining") != "1" {
		t.Fatalf("expected remaining header 1, got %s", w.Header().Get("X-RateLimit-Remaining"))
	}

	if w.Header().Get("X-RateLimit-Reset") == "" {
		t.Fatal("expected reset header to be set")
	}
}

func TestRateLimitPublicEndpointReturns429AfterLimitReached(t *testing.T) {
	router := setupRateLimitRouter(testRateLimitConfig())

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header")
	}

	if w.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Fatalf("expected remaining header 0, got %s", w.Header().Get("X-RateLimit-Remaining"))
	}

	if !strings.Contains(w.Body.String(), `"code":"rate_limited"`) {
		t.Fatalf("expected rate_limited error code, got: %s", w.Body.String())
	}
}

func TestRateLimitLoginEndpointUsesLoginLimit(t *testing.T) {
	router := setupRateLimitRouter(testRateLimitConfig())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("X-RateLimit-Limit") != "1" {
		t.Fatalf("expected login limit header 1, got %s", w.Header().Get("X-RateLimit-Limit"))
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRateLimitAuthenticatedEndpointUsesAuthLimit(t *testing.T) {
	router := setupRateLimitRouter(testRateLimitConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Header().Get("X-RateLimit-Limit") != "3" {
		t.Fatalf("expected auth limit header 3, got %s", w.Header().Get("X-RateLimit-Limit"))
	}
}

func TestLoadRateLimitConfigFromEnv(t *testing.T) {
	t.Setenv("RATE_LIMIT_PUBLIC", "10/1m")
	t.Setenv("RATE_LIMIT_AUTH", "20/2m")
	t.Setenv("RATE_LIMIT_LOGIN", "3/30s")

	config := LoadRateLimitConfigFromEnv()

	if config.Public.Limit != 10 {
		t.Fatalf("expected public limit 10, got %d", config.Public.Limit)
	}

	if config.Public.Window != time.Minute {
		t.Fatalf("expected public window 1m, got %v", config.Public.Window)
	}

	if config.Auth.Limit != 20 {
		t.Fatalf("expected auth limit 20, got %d", config.Auth.Limit)
	}

	if config.Auth.Window != 2*time.Minute {
		t.Fatalf("expected auth window 2m, got %v", config.Auth.Window)
	}

	if config.Login.Limit != 3 {
		t.Fatalf("expected login limit 3, got %d", config.Login.Limit)
	}

	if config.Login.Window != 30*time.Second {
		t.Fatalf("expected login window 30s, got %v", config.Login.Window)
	}
}
