package handlers

import "github.com/kodokbakar/go-ecommerce-api/internal/services"

type ErrorResponse struct {
	Error   string `json:"error" example:"bad_request"`
	Message string `json:"message" example:"invalid request body"`
}

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Go E-Commerce API is running"`
}

type AuthSuccessResponse struct {
	Message string                 `json:"message" example:"login successful"`
	Data    *services.AuthResponse `json:"data"`
}

type MeResponse struct {
	UserID string `json:"user_id" example:"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21"`
	Email  string `json:"email" example:"user@example.com"`
	Role   string `json:"role" example:"customer"`
}

type MessageResponse struct {
	Message string `json:"message" example:"admin route access granted"`
}
