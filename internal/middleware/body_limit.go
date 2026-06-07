package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

const (
	defaultBodyLimitAuth   int64 = 1024
	defaultBodyLimitAPI    int64 = 1048576
	defaultBodyLimitUpload int64 = 10485760
)

type BodyLimitConfig struct {
	Auth   int64
	API    int64
	Upload int64
}

func LoadBodyLimitConfigFromEnv() BodyLimitConfig {
	return BodyLimitConfig{
		Auth:   parseBodyLimitOrDefault(getEnvOrDefault("BODY_LIMIT_AUTH", strconv.FormatInt(defaultBodyLimitAuth, 10)), defaultBodyLimitAuth),
		API:    parseBodyLimitOrDefault(getEnvOrDefault("BODY_LIMIT_API", strconv.FormatInt(defaultBodyLimitAPI, 10)), defaultBodyLimitAPI),
		Upload: parseBodyLimitOrDefault(getEnvOrDefault("BODY_LIMIT_UPLOAD", strconv.FormatInt(defaultBodyLimitUpload, 10)), defaultBodyLimitUpload),
	}
}

func BodySizeLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit <= 0 {
			c.Next()
			return
		}

		if c.Request.Body == nil {
			c.Next()
			return
		}

		if c.Request.ContentLength > limit {
			respondBodyTooLarge(c, limit)
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)

		c.Next()
	}
}

func respondBodyTooLarge(c *gin.Context, limit int64) {
	response.PayloadTooLarge(
		c,
		"request body too large",
		fmt.Sprintf("request body exceeds %s limit", response.HumanReadableBytes(limit)),
	)
}

func parseBodyLimitOrDefault(value string, defaultValue int64) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	if parsed <= 0 {
		return defaultValue
	}

	return parsed
}
