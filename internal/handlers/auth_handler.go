package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
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
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	authResponse, err := h.authService.Register(c.Request.Context(), services.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", authResponse)
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
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	authResponse, err := h.authService.Login(c.Request.Context(), services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", authResponse)
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		response.BadRequest(c, err.Error(), nil)

	case errors.Is(err, services.ErrEmailAlreadyRegistered):
		response.Conflict(c, "email already registered", nil)

	case errors.Is(err, services.ErrInvalidCredentials):
		response.Unauthorized(c, "invalid email or password", nil)

	case errors.Is(err, services.ErrInactiveUser):
		response.Forbidden(c, "user account is inactive", nil)

	default:
		response.InternalServerError(c, "something went wrong", nil)
	}
}
