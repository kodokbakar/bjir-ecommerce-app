package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

var (
	ErrInvalidInput           = errors.New("invalid input")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInactiveUser           = errors.New("user account is inactive")
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	userRepository UserRepository
	jwtManager     *auth.JWTManager
}

func NewAuthService(userRepository UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := input.Password

	if err := validateRegisterInput(name, email, password); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepository.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.RoleCustomer,
		IsActive:     true,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			return nil, ErrEmailAlreadyRegistered
		}

		return nil, err
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return s.buildAuthResponse(user, token), nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := input.Password

	if err := validateLoginInput(email, password); err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return s.buildAuthResponse(user, token), nil
}

func (s *AuthService) buildAuthResponse(user *models.User, token string) *AuthResponse {
	return &AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.jwtManager.ExpiresIn().Seconds()),
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}
}

func validateRegisterInput(name string, email string, password string) error {
	if len(name) < 2 || len(name) > 100 {
		return fmt.Errorf("%w: name must be between 2 and 100 characters", ErrInvalidInput)
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%w: email is invalid", ErrInvalidInput)
	}

	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidInput)
	}

	if len([]byte(password)) > 72 {
		return fmt.Errorf("%w: password must not be longer than 72 bytes", ErrInvalidInput)
	}

	return nil
}

func validateLoginInput(email string, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%w: email is invalid", ErrInvalidInput)
	}

	if password == "" {
		return fmt.Errorf("%w: password is required", ErrInvalidInput)
	}

	return nil
}
