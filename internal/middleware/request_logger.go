package middleware

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const healthCheckPath = "/health"

type RequestLoggerConfig struct {
	Logger    *slog.Logger
	SkipPaths []string
}

func RequestLogger() gin.HandlerFunc {
	return RequestLoggerWithConfig(RequestLoggerConfig{
		Logger:    NewRequestLogLogger(os.Stdout),
		SkipPaths: []string{healthCheckPath},
	})
}

func RequestLoggerWithLogger(logger *slog.Logger) gin.HandlerFunc {
	return RequestLoggerWithConfig(RequestLoggerConfig{
		Logger:    logger,
		SkipPaths: []string{healthCheckPath},
	})
}

func RequestLoggerWithConfig(config RequestLoggerConfig) gin.HandlerFunc {
	logger := config.Logger
	if logger == nil {
		logger = NewRequestLogLogger(os.Stdout)
	}

	skipPaths := make(map[string]struct{}, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		if path == "" {
			continue
		}

		skipPaths[path] = struct{}{}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if _, shouldSkip := skipPaths[path]; shouldSkip {
			c.Next()
			return
		}

		startTime := time.Now()

		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		statusCode := c.Writer.Status()
		latency := time.Since(startTime)

		level := requestLogLevel(statusCode)
		levelLabel := requestLogLevelLabel(level)

		message := fmt.Sprintf(
			"[%s] %s | %s %s | %d | %s | %s",
			levelLabel,
			startTime.Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			latency.Round(time.Millisecond),
			clientIP,
		)

		logger.LogAttrs(
			context.Background(),
			level,
			message,
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status_code", statusCode),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
		)
	}
}

func NewRequestLogLogger(writer io.Writer) *slog.Logger {
	if writer == nil {
		writer = os.Stdout
	}

	return slog.New(&requestLogHandler{
		writer: writer,
		level:  slog.LevelInfo,
		mu:     &sync.Mutex{},
	})
}

type requestLogHandler struct {
	writer io.Writer
	level  slog.Level
	mu     *sync.Mutex
}

func (h *requestLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *requestLogHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := fmt.Fprintln(h.writer, record.Message)
	return err
}

func (h *requestLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *requestLogHandler) WithGroup(name string) slog.Handler {
	return h
}

func requestLogLevel(statusCode int) slog.Level {
	switch {
	case statusCode >= http.StatusInternalServerError:
		return slog.LevelError
	case statusCode >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func requestLogLevelLabel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARN"
	default:
		return "INFO"
	}
}
