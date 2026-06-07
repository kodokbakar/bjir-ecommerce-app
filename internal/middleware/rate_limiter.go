package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

const (
	defaultPublicRateLimit = "60/1m"
	defaultAuthRateLimit   = "120/1m"
	defaultLoginRateLimit  = "5/1m"

	rateLimitScopePublic = "public"
	rateLimitScopeAuth   = "auth"
	rateLimitScopeLogin  = "login"
)

type RateLimitRule struct {
	Limit  int
	Window time.Duration
}

type RateLimitConfig struct {
	Public RateLimitRule
	Auth   RateLimitRule
	Login  RateLimitRule
}

type rateLimitBucket struct {
	Count   int
	ResetAt time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	buckets     map[string]*rateLimitBucket
	config      RateLimitConfig
	lastCleanup time.Time
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		buckets:     make(map[string]*rateLimitBucket),
		config:      config,
		lastCleanup: time.Now(),
	}
}

func NewRateLimiterFromEnv() *RateLimiter {
	return NewRateLimiter(LoadRateLimitConfigFromEnv())
}

func LoadRateLimitConfigFromEnv() RateLimitConfig {
	return RateLimitConfig{
		Public: parseRateLimitOrDefault(getEnvOrDefault("RATE_LIMIT_PUBLIC", defaultPublicRateLimit), mustParseRateLimit(defaultPublicRateLimit)),
		Auth:   parseRateLimitOrDefault(getEnvOrDefault("RATE_LIMIT_AUTH", defaultAuthRateLimit), mustParseRateLimit(defaultAuthRateLimit)),
		Login:  parseRateLimitOrDefault(getEnvOrDefault("RATE_LIMIT_LOGIN", defaultLoginRateLimit), mustParseRateLimit(defaultLoginRateLimit)),
	}
}

func RateLimit() gin.HandlerFunc {
	return NewRateLimiterFromEnv().Middleware()
}

func RateLimitWithLimiter(limiter *RateLimiter) gin.HandlerFunc {
	return limiter.Middleware()
}

func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rule, scope := selectRateLimitRule(c.Request, l.config)
		key := buildRateLimitKey(scope, c.ClientIP())

		limit, remaining, resetAt, allowed := l.allow(key, rule)

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		if !allowed {
			retryAfterSeconds := secondsUntil(resetAt)

			c.Header("Retry-After", strconv.Itoa(retryAfterSeconds))

			message := fmt.Sprintf("too many requests, retry after %ds", retryAfterSeconds)
			response.RateLimited(c, "rate limit exceeded", message)
			return
		}

		c.Next()
	}
}

func (l *RateLimiter) allow(key string, rule RateLimitRule) (int, int, time.Time, bool) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.lastCleanup) > time.Minute {
		l.cleanupExpiredBucketsLocked(now)
		l.lastCleanup = now
	}

	bucket, exists := l.buckets[key]
	if !exists || !now.Before(bucket.ResetAt) {
		bucket = &rateLimitBucket{
			Count:   0,
			ResetAt: now.Add(rule.Window),
		}
		l.buckets[key] = bucket
	}

	if bucket.Count >= rule.Limit {
		return rule.Limit, 0, bucket.ResetAt, false
	}

	bucket.Count++

	remaining := rule.Limit - bucket.Count
	if remaining < 0 {
		remaining = 0
	}

	return rule.Limit, remaining, bucket.ResetAt, true
}

func (l *RateLimiter) cleanupExpiredBucketsLocked(now time.Time) {
	for key, bucket := range l.buckets {
		if !now.Before(bucket.ResetAt) {
			delete(l.buckets, key)
		}
	}
}

func selectRateLimitRule(req *http.Request, config RateLimitConfig) (RateLimitRule, string) {
	path := req.URL.Path
	method := req.Method

	if path == "/api/v1/auth/login" {
		return config.Login, rateLimitScopeLogin
	}

	if isAuthenticatedEndpoint(method, path) {
		return config.Auth, rateLimitScopeAuth
	}

	return config.Public, rateLimitScopePublic
}

func isAuthenticatedEndpoint(method string, path string) bool {
	if strings.HasPrefix(path, "/api/v1/me") ||
		strings.HasPrefix(path, "/api/v1/cart") ||
		strings.HasPrefix(path, "/api/v1/orders") ||
		strings.HasPrefix(path, "/api/v1/payments") ||
		strings.HasPrefix(path, "/api/v1/admin") {
		return true
	}

	if strings.HasPrefix(path, "/api/v1/categories") {
		return method == http.MethodPost ||
			method == http.MethodPut ||
			method == http.MethodPatch ||
			method == http.MethodDelete
	}

	if strings.HasPrefix(path, "/api/v1/products") {
		return method == http.MethodPost ||
			method == http.MethodPut ||
			method == http.MethodPatch ||
			method == http.MethodDelete
	}

	return false
}

func buildRateLimitKey(scope string, clientIP string) string {
	return scope + ":" + clientIP
}

func parseRateLimitOrDefault(value string, defaultRule RateLimitRule) RateLimitRule {
	rule, err := parseRateLimit(value)
	if err != nil {
		return defaultRule
	}

	return rule
}

func mustParseRateLimit(value string) RateLimitRule {
	rule, err := parseRateLimit(value)
	if err != nil {
		panic(err)
	}

	return rule
}

func parseRateLimit(value string) (RateLimitRule, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return RateLimitRule{}, fmt.Errorf("invalid rate limit format")
	}

	limit, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return RateLimitRule{}, err
	}

	if limit <= 0 {
		return RateLimitRule{}, fmt.Errorf("limit must be greater than 0")
	}

	window, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return RateLimitRule{}, err
	}

	if window <= 0 {
		return RateLimitRule{}, fmt.Errorf("window must be greater than 0")
	}

	return RateLimitRule{
		Limit:  limit,
		Window: window,
	}, nil
}

func secondsUntil(resetAt time.Time) int {
	seconds := int(time.Until(resetAt).Seconds())
	if seconds < 1 {
		return 1
	}

	return seconds
}
