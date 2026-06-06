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

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type payOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

// PayOrder godoc
// @Summary Pay order
// @Description Mock payment endpoint for customer order
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body payOrderRequest true "Payment request"
// @Success 201 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/pay [post]
func (h *PaymentHandler) PayOrder(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	var req payOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	payment, err := h.paymentService.PayOrder(c.Request.Context(), services.PayOrderInput{
		UserID:  userID,
		OrderID: req.OrderID,
		Method:  req.Method,
	})
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "payment successful", payment)
}

func handlePaymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidPaymentInput):
		response.BadRequest(c, "invalid payment input", err.Error())
	case errors.Is(err, models.ErrOrderNotFound):
		response.NotFound(c, "order not found", err.Error())
	case errors.Is(err, models.ErrOrderNotPayable):
		response.Conflict(c, "order is not payable", err.Error())
	case errors.Is(err, models.ErrPaymentAlreadyExists):
		response.Conflict(c, "payment already exists", err.Error())
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
}
