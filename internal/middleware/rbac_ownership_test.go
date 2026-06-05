package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

func newOwnershipTestJWTManager() *auth.JWTManager {
	return auth.NewJWTManager(config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})
}

func setupOwnershipTestRouter(jwtManager *auth.JWTManager) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	protected := r.Group("")
	protected.Use(AuthMiddleware(jwtManager))
	protected.GET("/users/:user_id/resource", RequireSelfOrAdmin("user_id"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})

	return r
}

func TestRequireSelfOrAdmin_NoToken_ReturnsUnauthorized(t *testing.T) {
	jwtManager := newOwnershipTestJWTManager()
	r := setupOwnershipTestRouter(jwtManager)

	req := httptest.NewRequest(http.MethodGet, "/users/user-id/resource", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRequireSelfOrAdmin_CustomerAccessOwnResource_ReturnsOK(t *testing.T) {
	jwtManager := newOwnershipTestJWTManager()
	r := setupOwnershipTestRouter(jwtManager)

	token, err := jwtManager.GenerateToken("user-id", "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/user-id/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRequireSelfOrAdmin_CustomerAccessOtherUserResource_ReturnsForbidden(t *testing.T) {
	jwtManager := newOwnershipTestJWTManager()
	r := setupOwnershipTestRouter(jwtManager)

	token, err := jwtManager.GenerateToken("user-id", "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/other-user-id/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRequireSelfOrAdmin_AdminAccessOtherUserResource_ReturnsOK(t *testing.T) {
	jwtManager := newOwnershipTestJWTManager()
	r := setupOwnershipTestRouter(jwtManager)

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/other-user-id/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}
