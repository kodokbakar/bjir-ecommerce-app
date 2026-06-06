package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

func TestRecovery_ReturnsJSONInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Recovery())

	r.GET("/panic", func(c *gin.Context) {
		panic("something went wrong")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}

	var body response.Body
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v. body: %s", err, w.Body.String())
	}

	if body.Success {
		t.Fatal("expected success false")
	}

	if body.Error == nil {
		t.Fatal("expected error body")
	}

	if body.Error.Code != response.CodeInternalServerError {
		t.Fatalf("expected code %s, got %s", response.CodeInternalServerError, body.Error.Code)
	}

	if body.Error.Message != "internal server error" {
		t.Fatalf("expected error message internal server error, got %s", body.Error.Message)
	}
}
