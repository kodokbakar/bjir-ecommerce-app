package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultCORSAllowOrigins = "*"
	defaultCORSAllowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	defaultCORSAllowHeaders = "Authorization,Content-Type"
	defaultCORSMaxAge       = 12 * time.Hour
)

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	MaxAge       time.Duration
}

func CORS() gin.HandlerFunc {
	return CORSWithConfig(LoadCORSConfigFromEnv())
}

func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	allowMethods := strings.Join(config.AllowMethods, ", ")
	allowHeaders := strings.Join(config.AllowHeaders, ", ")
	maxAgeSeconds := int(config.MaxAge.Seconds())

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if allowOrigin, ok := resolveAllowedOrigin(origin, config.AllowOrigins); ok {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Methods", allowMethods)
			c.Header("Access-Control-Allow-Headers", allowHeaders)
			c.Header("Access-Control-Max-Age", strconv.Itoa(maxAgeSeconds))

			if allowOrigin != "*" {
				c.Header("Vary", "Origin")
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func LoadCORSConfigFromEnv() CORSConfig {
	return CORSConfig{
		AllowOrigins: splitCSV(getEnvOrDefault("CORS_ALLOW_ORIGINS", defaultCORSAllowOrigins)),
		AllowMethods: splitCSV(getEnvOrDefault("CORS_ALLOW_METHODS", defaultCORSAllowMethods)),
		AllowHeaders: splitCSV(getEnvOrDefault("CORS_ALLOW_HEADERS", defaultCORSAllowHeaders)),
		MaxAge:       getDurationEnvOrDefault("CORS_MAX_AGE", defaultCORSMaxAge),
	}
}

func resolveAllowedOrigin(origin string, allowOrigins []string) (string, bool) {
	for _, allowedOrigin := range allowOrigins {
		if allowedOrigin == "*" {
			return "*", true
		}

		if origin != "" && origin == allowedOrigin {
			return origin, true
		}
	}

	return "", false
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}
