package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
)

func TestAuthMiddlewareWithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	token, err := jwtManager.GenerateToken(
		"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		"user@example.com",
		"customer",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	r := gin.New()
	r.GET("/protected", middleware.AuthMiddleware(jwtManager), func(c *gin.Context) {
		userID, ok := middleware.GetCurrentUserID(c)
		if !ok {
			t.Fatal("expected user id in context")
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareWithoutAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	r := gin.New()
	r.GET("/protected", middleware.AuthMiddleware(jwtManager), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "should not reach this handler",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareWithInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	r := gin.New()
	r.GET("/protected", middleware.AuthMiddleware(jwtManager), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "should not reach this handler",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRequireRoleWithAllowedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/admin",
		func(c *gin.Context) {
			c.Set(middleware.ContextUserRole, "admin")
			c.Next()
		},
		middleware.RequireRole("admin"),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRequireRoleWithForbiddenRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/admin",
		func(c *gin.Context) {
			c.Set(middleware.ContextUserRole, "customer")
			c.Next()
		},
		middleware.RequireRole("admin"),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "should not reach this handler",
			})
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d. body: %s", w.Code, w.Body.String())
	}
}
