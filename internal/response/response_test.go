package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/success", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "request successful", gin.H{
			"id": "123",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/success", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body response.Body
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !body.Success {
		t.Fatal("expected success true")
	}

	if body.Message != "request successful" {
		t.Fatalf("expected message request successful, got %s", body.Message)
	}

	if body.Error != nil {
		t.Fatal("expected error nil")
	}
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/error", func(c *gin.Context) {
		response.BadRequest(c, "invalid request body", "email is required")
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var body response.Body
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Success {
		t.Fatal("expected success false")
	}

	if body.Message != "invalid request body" {
		t.Fatalf("expected message invalid request body, got %s", body.Message)
	}

	if body.Error == nil {
		t.Fatal("expected error body, got nil")
	}

	if body.Error.Code != response.CodeBadRequest {
		t.Fatalf("expected code bad_request, got %s", body.Error.Code)
	}
}
