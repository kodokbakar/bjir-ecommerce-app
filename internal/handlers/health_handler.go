package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check whether the API server is running.
// @Tags System
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "Go E-Commerce API is running", gin.H{
		"status": "ok",
	})
}
