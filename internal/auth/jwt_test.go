package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

func TestGenerateAndValidateToken(t *testing.T) {
	manager := NewJWTManager(config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	token, err := manager.GenerateToken(
		"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		"user@example.com",
		"customer",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected token to be valid, got %v", err)
	}

	if claims.UserID != "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21" {
		t.Fatalf("expected user id to match, got %s", claims.UserID)
	}

	if claims.Email != "user@example.com" {
		t.Fatalf("expected email to match, got %s", claims.Email)
	}

	if claims.Role != "customer" {
		t.Fatalf("expected role to match, got %s", claims.Role)
	}
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
	manager := NewJWTManager(config.JWTConfig{
		Secret:    "correct_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	token, err := manager.GenerateToken(
		"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		"user@example.com",
		"customer",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrongManager := NewJWTManager(config.JWTConfig{
		Secret:    "wrong_secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	_, err = wrongManager.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}

	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateExpiredToken(t *testing.T) {
	manager := NewJWTManager(config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: -time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	token, err := manager.GenerateToken(
		"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		"user@example.com",
		"customer",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}

	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestExtractBearerToken(t *testing.T) {
	token, err := ExtractBearerToken("Bearer abc.def.ghi")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "abc.def.ghi" {
		t.Fatalf("expected token to match, got %s", token)
	}
}

func TestExtractBearerTokenInvalidHeader(t *testing.T) {
	_, err := ExtractBearerToken("abc.def.ghi")
	if err == nil {
		t.Fatal("expected error for invalid authorization header")
	}

	if !errors.Is(err, ErrInvalidAuthHeader) {
		t.Fatalf("expected ErrInvalidAuthHeader, got %v", err)
	}
}
