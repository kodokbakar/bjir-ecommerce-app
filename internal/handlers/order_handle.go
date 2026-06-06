package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/response"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	order, err := h.orderService.Checkout(c.Request.Context(), userID)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "checkout successful", order)
}

func handleOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidOrderInput):
		response.BadRequest(c, "invalid order input", err.Error())
	case errors.Is(err, models.ErrCartEmpty):
		response.BadRequest(c, "cart is empty", nil)
	case errors.Is(err, models.ErrInsufficientStock):
		response.Conflict(c, "insufficient product stock", nil)
	case errors.Is(err, models.ErrProductNotFound):
		response.NotFound(c, "product not found", nil)
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
}
