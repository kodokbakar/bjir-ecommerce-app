package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

var shuttingDown atomic.Bool

func SetShuttingDown(value bool) {
	shuttingDown.Store(value)
}

func IsShuttingDown() bool {
	return shuttingDown.Load()
}

// HealthCheck godoc
// @Summary Health check
// @Description Check API health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} ErrorResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	if IsShuttingDown() {
		response.Error(
			c,
			http.StatusServiceUnavailable,
			response.CodeServiceUnavailable,
			"server is shutting down",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, "Go E-Commerce API is running", HealthData{
		Status: "ok",
	})
}
