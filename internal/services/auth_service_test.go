package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type mockUserRepository struct {
	usersByEmail map[string]*models.User
	createErr    error
	getErr       error
	createdUsers []*models.User
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
	m.createdUsers = append(m.createdUsers, user)

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

func TestAuthServiceRegisterValidData(t *testing.T) {
	repo := newMockUserRepository()
	service := services.NewAuthService(repo, newTestJWTManager())

	response, err := service.Register(context.Background(), services.RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.AccessToken == "" {
		t.Fatal("expected access token, got empty string")
	}

	if response.TokenType != "Bearer" {
		t.Fatalf("expected token type Bearer, got %s", response.TokenType)
	}

	if response.User.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", response.User.Email)
	}

	if response.User.Role != models.RoleCustomer {
		t.Fatalf("expected role customer, got %s", response.User.Role)
	}

	if len(repo.createdUsers) != 1 {
		t.Fatalf("expected 1 created user, got %d", len(repo.createdUsers))
	}

	if repo.createdUsers[0].PasswordHash == "password123" {
		t.Fatal("expected password to be hashed, got plain password")
	}
}

func TestAuthServiceRegisterDuplicateEmail(t *testing.T) {
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

	service := services.NewAuthService(repo, newTestJWTManager())

	_, err := service.Register(context.Background(), services.RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrEmailAlreadyRegistered) {
		t.Fatalf("expected ErrEmailAlreadyRegistered, got %v", err)
	}
}

func TestAuthServiceRegisterWeakPassword(t *testing.T) {
	repo := newMockUserRepository()
	service := services.NewAuthService(repo, newTestJWTManager())

	_, err := service.Register(context.Background(), services.RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAuthServiceLoginValidCredentials(t *testing.T) {
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

	service := services.NewAuthService(repo, newTestJWTManager())

	response, err := service.Login(context.Background(), services.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.AccessToken == "" {
		t.Fatal("expected access token, got empty string")
	}

	if response.User.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", response.User.Email)
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
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

	service := services.NewAuthService(repo, newTestJWTManager())

	_, err = service.Login(context.Background(), services.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthServiceLoginNonexistentUser(t *testing.T) {
	repo := newMockUserRepository()
	service := services.NewAuthService(repo, newTestJWTManager())

	_, err := service.Login(context.Background(), services.LoginInput{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthServiceLoginInactiveUser(t *testing.T) {
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

	service := services.NewAuthService(repo, newTestJWTManager())

	_, err = service.Login(context.Background(), services.LoginInput{
		Email:    "inactive@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrInactiveUser) {
		t.Fatalf("expected ErrInactiveUser, got %v", err)
	}
}
