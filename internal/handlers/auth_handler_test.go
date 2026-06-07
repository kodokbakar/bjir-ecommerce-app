package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type mockUserRepository struct {
	usersByEmail map[string]*models.User
	createErr    error
	getErr       error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		usersByEmail: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}

	if _, exists := m.usersByEmail[user.Email]; exists {
		return models.ErrDuplicateEmail
	}

	if user.ID == "" {
		user.ID = "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21"
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = user.CreatedAt
	}

	m.usersByEmail[user.Email] = user

	return nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}

	user, exists := m.usersByEmail[email]
	if !exists {
		return nil, models.ErrUserNotFound
	}

	return user, nil
}

func newTestJWTManager() *auth.JWTManager {
	return auth.NewJWTManager(config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})
}

func setupAuthRouter(repo *mockUserRepository) *gin.Engine {
	r, _ := setupAuthRouterWithJWT(repo)
	return r
}

func setupAuthRouterWithJWT(repo *mockUserRepository) (*gin.Engine, *auth.JWTManager) {
	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager()
	authService := services.NewAuthService(repo, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.New()
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	protected.GET("/me", handlers.Me)

	return r, jwtManager
}

func performRequest(r http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func performRequestWithToken(r http.Handler, method string, path string, body string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestRegisterHandlerValidData(t *testing.T) {
	repo := newMockUserRepository()
	r := setupAuthRouter(repo)

	body := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", w.Code, w.Body.String())
	}

	if _, exists := repo.usersByEmail["test@example.com"]; !exists {
		t.Fatal("expected user to be created")
	}
}

func TestRegisterHandlerDuplicateEmail(t *testing.T) {
	repo := newMockUserRepository()
	repo.usersByEmail["test@example.com"] = &models.User{
		ID:           "existing-user-id",
		Name:         "Existing User",
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
		Role:         models.RoleCustomer,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r := setupAuthRouter(repo)

	body := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRegisterHandlerUnexpectedError(t *testing.T) {
	repo := newMockUserRepository()
	repo.createErr = errors.New("database connection lost")

	r := setupAuthRouter(repo)

	body := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), false)
	assertErrorCode(t, w.Body.Bytes(), "internal_server_error")
}

func TestRegisterHandlerWeakPassword(t *testing.T) {
	repo := newMockUserRepository()
	r := setupAuthRouter(repo)

	body := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRegisterHandlerInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing name",
			body: `{
				"email": "test@example.com",
				"password": "password123"
			}`,
		},
		{
			name: "invalid email",
			body: `{
				"name": "Test User",
				"email": "not-an-email",
				"password": "password123"
			}`,
		},
		{
			name: "short password",
			body: `{
				"name": "Test User",
				"email": "test@example.com",
				"password": "123"
			}`,
		},
		{
			name: "empty body",
			body: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			r := setupAuthRouter(repo)

			w := performRequest(r, http.MethodPost, "/api/v1/auth/register", tt.body)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
			}

			assertStandardResponse(t, w.Body.Bytes(), false)
			assertErrorCode(t, w.Body.Bytes(), "bad_request")
		})
	}
}

func TestLoginHandlerValidCredentials(t *testing.T) {
	repo := newMockUserRepository()

	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.usersByEmail["test@example.com"] = &models.User{
		ID:           "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Role:         models.RoleCustomer,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r := setupAuthRouter(repo)

	body := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestLoginHandlerWrongPassword(t *testing.T) {
	repo := newMockUserRepository()

	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.usersByEmail["test@example.com"] = &models.User{
		ID:           "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Role:         models.RoleCustomer,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r := setupAuthRouter(repo)

	body := `{
		"email": "test@example.com",
		"password": "wrongpassword"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestLoginHandlerNonexistentUser(t *testing.T) {
	repo := newMockUserRepository()
	r := setupAuthRouter(repo)

	body := `{
		"email": "notfound@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestLoginHandlerUnexpectedError(t *testing.T) {
	repo := newMockUserRepository()
	repo.getErr = errors.New("database connection lost")

	r := setupAuthRouter(repo)

	body := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), false)
	assertErrorCode(t, w.Body.Bytes(), "internal_server_error")
}

func TestLoginHandlerInactiveUser(t *testing.T) {
	repo := newMockUserRepository()

	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.usersByEmail["inactive@example.com"] = &models.User{
		ID:           "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		Name:         "Inactive User",
		Email:        "inactive@example.com",
		PasswordHash: passwordHash,
		Role:         models.RoleCustomer,
		IsActive:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r := setupAuthRouter(repo)

	body := `{
		"email": "inactive@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestLoginHandlerInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "invalid email",
			body: `{
				"email": "not-an-email",
				"password": "password123"
			}`,
		},
		{
			name: "missing password",
			body: `{
				"email": "test@example.com"
			}`,
		},
		{
			name: "empty body",
			body: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			r := setupAuthRouter(repo)

			w := performRequest(r, http.MethodPost, "/api/v1/auth/login", tt.body)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
			}

			assertStandardResponse(t, w.Body.Bytes(), false)
			assertErrorCode(t, w.Body.Bytes(), "bad_request")
		})
	}
}

func assertStandardResponse(t *testing.T, body []byte, expectedSuccess bool) {
	t.Helper()

	var responseBody struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    any    `json:"data"`
		Error   any    `json:"error"`
	}

	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("failed to decode response body: %v. body: %s", err, string(body))
	}

	if responseBody.Success != expectedSuccess {
		t.Fatalf("expected success %v, got %v. body: %s", expectedSuccess, responseBody.Success, string(body))
	}

	if responseBody.Message == "" {
		t.Fatalf("expected message field, got empty. body: %s", string(body))
	}
}

func assertErrorCode(t *testing.T, body []byte, expectedCode string) {
	t.Helper()

	var responseBody struct {
		Success bool `json:"success"`
		Error   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("failed to decode response body: %v. body: %s", err, string(body))
	}

	if responseBody.Success {
		t.Fatalf("expected success false, got true. body: %s", string(body))
	}

	if responseBody.Error.Code != expectedCode {
		t.Fatalf("expected error code %s, got %s. body: %s", expectedCode, responseBody.Error.Code, string(body))
	}

	if responseBody.Error.Message == "" {
		t.Fatalf("expected error message, got empty. body: %s", string(body))
	}
}

func TestMeHandlerSuccess(t *testing.T) {
	repo := newMockUserRepository()
	r, jwtManager := setupAuthRouterWithJWT(repo)

	token, err := jwtManager.GenerateToken(
		"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		"test@example.com",
		models.RoleCustomer,
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := performRequestWithToken(r, http.MethodGet, "/api/v1/me", "", token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), true)

	var responseBody struct {
		Success bool `json:"success"`
		Data    struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
			Role   string `json:"role"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("failed to decode response body: %v. body: %s", err, w.Body.String())
	}

	if responseBody.Data.UserID != "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21" {
		t.Fatalf("expected user_id from token, got %s", responseBody.Data.UserID)
	}

	if responseBody.Data.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", responseBody.Data.Email)
	}

	if responseBody.Data.Role != models.RoleCustomer {
		t.Fatalf("expected role customer, got %s", responseBody.Data.Role)
	}
}

func TestMeHandlerWithoutToken(t *testing.T) {
	repo := newMockUserRepository()
	r, _ := setupAuthRouterWithJWT(repo)

	w := performRequest(r, http.MethodGet, "/api/v1/me", "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), false)
	assertErrorCode(t, w.Body.Bytes(), "unauthorized")
}

func TestMeHandlerInvalidToken(t *testing.T) {
	repo := newMockUserRepository()
	r, _ := setupAuthRouterWithJWT(repo)

	w := performRequestWithToken(r, http.MethodGet, "/api/v1/me", "", "invalid-token")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), false)
	assertErrorCode(t, w.Body.Bytes(), "unauthorized")
}

func TestRegisterHandlerBodyTooLarge(t *testing.T) {
	repo := newMockUserRepository()

	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager()
	authService := services.NewAuthService(repo, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.New()
	r.Use(middleware.BodySizeLimit(20))
	r.POST("/api/v1/auth/register", authHandler.Register)

	body := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123"
	}`

	w := performRequest(r, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), false)
	assertErrorCode(t, w.Body.Bytes(), "payload_too_large")
}

func TestRegisterHandlerBodyTooLargeWithoutContentLength(t *testing.T) {
	repo := newMockUserRepository()

	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager()
	authService := services.NewAuthService(repo, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.New()
	r.Use(middleware.BodySizeLimit(20))
	r.POST("/api/v1/auth/register", authHandler.Register)

	body := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = -1

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d. body: %s", w.Code, w.Body.String())
	}

	assertStandardResponse(t, w.Body.Bytes(), false)
	assertErrorCode(t, w.Body.Bytes(), "payload_too_large")
}
