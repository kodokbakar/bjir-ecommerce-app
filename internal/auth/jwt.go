package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

var (
	ErrEmptyToken        = errors.New("token is empty")
	ErrInvalidToken      = errors.New("token is invalid")
	ErrInvalidAuthHeader = errors.New("authorization header must use Bearer token")
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

type JWTManager struct {
	secret    []byte
	expiresIn time.Duration
	issuer    string
}

func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:    []byte(cfg.Secret),
		expiresIn: cfg.ExpiresIn,
		issuer:    cfg.Issuer,
	}
}

func (m *JWTManager) ExpiresIn() time.Duration {
	return m.expiresIn
}

func (m *JWTManager) GenerateToken(userID string, email string, role string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("userID is required")
	}

	if email == "" {
		return "", fmt.Errorf("email is required")
	}

	if role == "" {
		return "", fmt.Errorf("role is required")
	}

	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrEmptyToken
	}

	claims := &Claims{}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
	)

	token, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	parsedClaims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if parsedClaims.UserID == "" {
		return nil, fmt.Errorf("%w: user_id is empty", ErrInvalidToken)
	}

	if parsedClaims.Subject != parsedClaims.UserID {
		return nil, fmt.Errorf("%w: subject does not match user_id", ErrInvalidToken)
	}

	return parsedClaims, nil
}

func ExtractBearerToken(authorizationHeader string) (string, error) {
	authorizationHeader = strings.TrimSpace(authorizationHeader)
	if authorizationHeader == "" {
		return "", ErrInvalidAuthHeader
	}

	parts := strings.Fields(authorizationHeader)
	if len(parts) != 2 {
		return "", ErrInvalidAuthHeader
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrInvalidAuthHeader
	}

	if strings.TrimSpace(parts[1]) == "" {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}
