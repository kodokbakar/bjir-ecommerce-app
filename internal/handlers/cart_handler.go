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

type CartHandler struct {
	cartService services.CartService
}

type addCartItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

type updateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *CartHandler) AddCartItem(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	var req addCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	item, err := h.cartService.AddItem(c.Request.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		handleCartError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "cart item added successfully", item)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		handleCartError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "cart retrieved successfully", cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	itemID := c.Param("id")

	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	item, err := h.cartService.UpdateItem(c.Request.Context(), userID, itemID, req.Quantity)
	if err != nil {
		handleCartError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "cart item updated successfully", item)
}

func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c, "unauthorized", "user id not found in context")
		return
	}

	itemID := c.Param("id")

	if err := h.cartService.DeleteItem(c.Request.Context(), userID, itemID); err != nil {
		handleCartError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleCartError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidCartInput):
		response.BadRequest(c, "invalid cart input", err.Error())
	case errors.Is(err, models.ErrProductNotFound):
		response.NotFound(c, "product not found", nil)
	case errors.Is(err, models.ErrCartItemNotFound):
		response.NotFound(c, "cart item not found", nil)
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
}
