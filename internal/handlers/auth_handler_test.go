package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
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
	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager()
	authService := services.NewAuthService(repo, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.New()
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)

	return r
}

func performRequest(r http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

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
