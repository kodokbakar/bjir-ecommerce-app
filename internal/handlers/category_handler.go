package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/response"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

type createCategoryRequest struct {
	ParentID    *string `json:"parent_id"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
}

type updateCategoryRequest struct {
	ParentID    *string `json:"parent_id"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req createCategoryRequest

	if !bindAndValidateJSON(c, &req) {
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), services.CreateCategoryInput{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		h.handleCategoryError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "category created successfully", category)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	pagination := GetPaginationQuery(c)

	result, err := h.categoryService.GetAll(c.Request.Context(), services.CategoryListInput{
		Page:  pagination.Page,
		Limit: pagination.Limit,
	})
	if err != nil {
		h.handleCategoryError(c, err)
		return
	}

	meta := gin.H{
		"page":        result.Page,
		"limit":       result.Limit,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	}

	response.SuccessWithMeta(c, http.StatusOK, "categories retrieved successfully", result.Categories, meta)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleCategoryError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "category retrieved successfully", category)
}

func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		h.handleCategoryError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "category retrieved successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req updateCategoryRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, services.UpdateCategoryInput{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		h.handleCategoryError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		h.handleCategoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) handleCategoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidCategoryInput):
		response.BadRequest(c, "invalid category input", err.Error())
	case errors.Is(err, models.ErrCategoryNotFound):
		response.NotFound(c, "category not found", err.Error())
	case errors.Is(err, models.ErrCategoryAlreadyExists):
		response.Conflict(c, "category already exists", err.Error())
	case errors.Is(err, models.ErrCategoryHasProducts):
		response.Conflict(c, "category cannot be deleted because it has related products", err.Error())
	case errors.Is(err, models.ErrCategoryHasChildren):
		response.Conflict(c, "category cannot be deleted because it has child categories", err.Error())
	default:
		response.InternalServerError(c, "internal server error", err.Error())
	}
}
