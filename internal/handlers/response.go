package handlers

import "github.com/kodokbakar/go-ecommerce-api/internal/services"

type ErrorResponse struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"invalid request body"`
	Error   ErrorResponseBody `json:"error"`
}

type ErrorResponseBody struct {
	Code    string `json:"code" example:"bad_request"`
	Details any    `json:"details,omitempty"`
}

type HealthResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"Go E-Commerce API is running"`
	Data    HealthData `json:"data"`
}

type HealthData struct {
	Status string `json:"status" example:"ok"`
}

type AuthSuccessResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"login successful"`
	Data    *services.AuthResponse `json:"data"`
}

type MeResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"current user retrieved successfully"`
	Data    MeData `json:"data"`
}

type MeData struct {
	UserID string `json:"user_id" example:"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21"`
	Email  string `json:"email" example:"user@example.com"`
	Role   string `json:"role" example:"customer"`
}

type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"admin route access granted"`
}
