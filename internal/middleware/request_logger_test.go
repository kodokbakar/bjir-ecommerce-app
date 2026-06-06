package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRequestLogger(buffer *bytes.Buffer) gin.HandlerFunc {
	return RequestLoggerWithConfig(RequestLoggerConfig{
		Logger:    NewRequestLogLogger(buffer),
		SkipPaths: []string{healthCheckPath},
	})
}

func TestRequestLogger_LogsSuccessfulRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buffer bytes.Buffer

	r := gin.New()
	r.Use(newTestRequestLogger(&buffer))

	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	output := buffer.String()

	pattern := regexp.MustCompile(`^\[INFO\] \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \| POST /api/v1/auth/login \| 200 \| .+ \| .+\n$`)
	if !pattern.MatchString(output) {
		t.Fatalf("expected log to match request log format, got: %s", output)
	}

	if strings.Contains(output, "msg=") {
		t.Fatalf("expected plain request log format without slog msg wrapper, got: %s", output)
	}
}

func TestRequestLogger_LogsWarnForClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buffer bytes.Buffer

	r := gin.New()
	r.Use(newTestRequestLogger(&buffer))

	r.GET("/bad-request", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
	})

	req := httptest.NewRequest(http.MethodGet, "/bad-request", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	output := buffer.String()

	if !strings.HasPrefix(output, "[WARN]") {
		t.Fatalf("expected WARN log, got: %s", output)
	}

	if !strings.Contains(output, "GET /bad-request | 400") {
		t.Fatalf("expected method, path, and status code in log, got: %s", output)
	}
}

func TestRequestLogger_LogsErrorForServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buffer bytes.Buffer

	r := gin.New()
	r.Use(newTestRequestLogger(&buffer))

	r.GET("/server-error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
	})

	req := httptest.NewRequest(http.MethodGet, "/server-error", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	output := buffer.String()

	if !strings.HasPrefix(output, "[ERROR]") {
		t.Fatalf("expected ERROR log, got: %s", output)
	}

	if !strings.Contains(output, "GET /server-error | 500") {
		t.Fatalf("expected method, path, and status code in log, got: %s", output)
	}
}

func TestRequestLogger_SkipsHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buffer bytes.Buffer

	r := gin.New()
	r.Use(newTestRequestLogger(&buffer))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if buffer.String() != "" {
		t.Fatalf("expected health check to be skipped, got log: %s", buffer.String())
	}
}

func TestRequestLogger_SkipsConfiguredPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buffer bytes.Buffer

	r := gin.New()
	r.Use(RequestLoggerWithConfig(RequestLoggerConfig{
		Logger:    NewRequestLogLogger(&buffer),
		SkipPaths: []string{"/metrics"},
	}))

	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if buffer.String() != "" {
		t.Fatalf("expected configured path to be skipped, got log: %s", buffer.String())
	}
}
