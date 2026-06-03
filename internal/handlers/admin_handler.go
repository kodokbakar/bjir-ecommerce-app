package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

// AdminPing godoc
// @Summary Admin ping
// @Description Access admin-only endpoint.
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/v1/admin/ping [get]
func AdminPing(c *gin.Context) {
	response.Success(c, http.StatusOK, "admin route access granted", nil)
}
