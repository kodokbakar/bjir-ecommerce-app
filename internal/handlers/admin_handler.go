package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{
		"message": "admin route access granted",
	})
}
