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

type updateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
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
	case errors.Is(err, models.ErrOrderNotFound):
		response.NotFound(c, "order not found", nil)
	case errors.Is(err, models.ErrInvalidOrderStatus):
		response.BadRequest(c, "invalid order status", err.Error())
	case errors.Is(err, models.ErrInvalidOrderStatusTransition):
		response.Conflict(c, "invalid order status transition", err.Error())
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	pagination := GetPaginationQuery(c)

	result, err := h.orderService.GetMyOrders(c.Request.Context(), userID, services.OrderListInput{
		Page:  pagination.Page,
		Limit: pagination.Limit,
	})
	if err != nil {
		handleOrderError(c, err)
		return
	}

	meta := gin.H{
		"page":        result.Page,
		"limit":       result.Limit,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	}

	response.SuccessWithMeta(c, http.StatusOK, "orders retrieved successfully", result.Orders, meta)
}

func (h *OrderHandler) GetMyOrderDetail(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	orderID := c.Param("id")

	order, err := h.orderService.GetMyOrderDetail(c.Request.Context(), userID, orderID)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order retrieved successfully", order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req updateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	order, err := h.orderService.UpdateStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order status updated successfully", order)
}
