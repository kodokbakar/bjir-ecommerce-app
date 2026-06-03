package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
)

const (
	ContextAuthClaims = "auth_claims"
	ContextUserID     = "user_id"
	ContextUserEmail  = "user_email"
	ContextUserRole   = "user_role"
)

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		tokenString, err := auth.ExtractBearerToken(authorizationHeader)
		if err != nil {
			respondUnauthorized(c, "missing or invalid authorization header")
			return
		}

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			respondUnauthorized(c, "invalid or expired token")
			return
		}

		c.Set(ContextAuthClaims, claims)
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserEmail, claims.Email)
		c.Set(ContextUserRole, claims.Role)

		c.Next()
	}
}

func GetAuthClaims(c *gin.Context) (*auth.Claims, bool) {
	value, exists := c.Get(ContextAuthClaims)
	if !exists {
		return nil, false
	}

	claims, ok := value.(*auth.Claims)
	return claims, ok
}

func GetCurrentUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(ContextUserID)
	if !exists {
		return "", false
	}

	userID, ok := value.(string)
	return userID, ok
}

func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	value, exists := c.Get(ContextUserEmail)
	if !exists {
		return "", false
	}

	email, ok := value.(string)
	return email, ok
}

func GetCurrentUserRole(c *gin.Context) (string, bool) {
	value, exists := c.Get(ContextUserRole)
	if !exists {
		return "", false
	}

	role, ok := value.(string)
	return role, ok
}

func respondUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error":   "unauthorized",
		"message": message,
	})
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRole, ok := GetCurrentUserRole(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "user role not found in context",
			})
			return
		}

		for _, role := range allowedRoles {
			if currentRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "you do not have permission to access this resource",
		})
	}
}

func IsUnauthorizedError(err error) bool {
	return errors.Is(err, auth.ErrEmptyToken) ||
		errors.Is(err, auth.ErrInvalidToken) ||
		errors.Is(err, auth.ErrInvalidAuthHeader)
}
