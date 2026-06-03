package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

// Me godoc
// @Summary Get current user from token
// @Description Return user identity extracted from validated Bearer token.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MeResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/me [get]
func Me(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	email, _ := middleware.GetCurrentUserEmail(c)
	role, _ := middleware.GetCurrentUserRole(c)

	response.Success(c, http.StatusOK, "current user retrieved successfully", gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}
