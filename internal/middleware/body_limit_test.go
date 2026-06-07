package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupBodyLimitRouter(limit int64) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BodySizeLimit(limit))

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	return router
}

func TestBodySizeLimitAllowsNormalRequest(t *testing.T) {
	router := setupBodyLimitRouter(10)

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("12345"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestBodySizeLimitReturns413WhenContentLengthExceedsLimit(t *testing.T) {
	router := setupBodyLimitRouter(10)

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("12345678901"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"code":"payload_too_large"`) {
		t.Fatalf("expected payload_too_large code, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"message":"request body too large"`) {
		t.Fatalf("expected request body too large message, got: %s", w.Body.String())
	}
}

func TestLoadBodyLimitConfigFromEnv(t *testing.T) {
	t.Setenv("BODY_LIMIT_AUTH", "2048")
	t.Setenv("BODY_LIMIT_API", "2097152")
	t.Setenv("BODY_LIMIT_UPLOAD", "5242880")

	config := LoadBodyLimitConfigFromEnv()

	if config.Auth != 2048 {
		t.Fatalf("expected auth limit 2048, got %d", config.Auth)
	}

	if config.API != 2097152 {
		t.Fatalf("expected api limit 2097152, got %d", config.API)
	}

	if config.Upload != 5242880 {
		t.Fatalf("expected upload limit 5242880, got %d", config.Upload)
	}
}

func TestLoadBodyLimitConfigFromEnvUsesDefaultForInvalidValue(t *testing.T) {
	t.Setenv("BODY_LIMIT_AUTH", "invalid")
	t.Setenv("BODY_LIMIT_API", "-1")
	t.Setenv("BODY_LIMIT_UPLOAD", "0")

	config := LoadBodyLimitConfigFromEnv()

	if config.Auth != defaultBodyLimitAuth {
		t.Fatalf("expected default auth limit %d, got %d", defaultBodyLimitAuth, config.Auth)
	}

	if config.API != defaultBodyLimitAPI {
		t.Fatalf("expected default api limit %d, got %d", defaultBodyLimitAPI, config.API)
	}

	if config.Upload != defaultBodyLimitUpload {
		t.Fatalf("expected default upload limit %d, got %d", defaultBodyLimitUpload, config.Upload)
	}
}
