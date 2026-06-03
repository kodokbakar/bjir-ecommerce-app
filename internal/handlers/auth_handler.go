package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100" example:"Test User"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Register godoc
// @Summary Register user
// @Description Create a new customer account and return an access token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request body"
// @Success 201 {object} AuthSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad_request",
			"message": "invalid request body",
		})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), services.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"data":    response,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password, then return an access token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request body"
// @Success 200 {object} AuthSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad_request",
			"message": "invalid request body",
		})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data":    response,
	})
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad_request",
			"message": err.Error(),
		})

	case errors.Is(err, services.ErrEmailAlreadyRegistered):
		c.JSON(http.StatusConflict, gin.H{
			"error":   "conflict",
			"message": "email already registered",
		})

	case errors.Is(err, services.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "invalid email or password",
		})

	case errors.Is(err, services.ErrInactiveUser):
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "user account is inactive",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_server_error",
			"message": "something went wrong",
		})
	}
}
