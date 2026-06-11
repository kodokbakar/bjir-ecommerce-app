package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/response"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type createProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required,max=150"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
}

type updateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required,max=150"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
}

type reorderProductImagesRequest struct {
	Images []reorderProductImageRequest `json:"images" binding:"required,min=1,dive"`
}

type reorderProductImageRequest struct {
	ID        string `json:"id" binding:"required"`
	SortOrder int    `json:"sort_order" binding:"gte=0"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req createProductRequest

	if !bindAndValidateJSON(c, &req) {
		return
	}

	product, err := h.productService.Create(c.Request.Context(), services.CreateProductInput{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "product created successfully", product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	categoryID := c.Query("category_id")
	categorySlug := c.Query("category")
	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	pagination := GetPaginationQuery(c)

	result, err := h.productService.GetAll(c.Request.Context(), services.ProductListInput{
		CategoryID:   categoryID,
		CategorySlug: categorySlug,
		Search:       search,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	})
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	meta := gin.H{
		"page":        result.Page,
		"limit":       result.Limit,
		"total":       result.Total,
		"total_pages": result.TotalPages,
		"sort_by":     result.SortBy,
		"sort_order":  result.SortOrder,
		"category_id": result.CategoryID,
		"category":    result.CategorySlug,
		"search":      result.Search,
	}

	response.SuccessWithMeta(c, http.StatusOK, "products retrieved successfully", result.Products, meta)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "product retrieved successfully", product)
}

func (h *ProductHandler) GetProductImages(c *gin.Context) {
	productID := c.Param("id")

	images, err := h.productService.GetImages(c.Request.Context(), productID)
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "product images retrieved successfully", images)
}

func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.productService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "product retrieved successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req updateProductRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, services.UpdateProductInput{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		h.handleProductError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) handleProductError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidProductInput):
		response.BadRequest(c, "invalid product input", err.Error())
	case errors.Is(err, models.ErrProductNotFound):
		response.NotFound(c, "product not found", err.Error())
	case errors.Is(err, models.ErrProductImageNotFound):
		response.NotFound(c, "product image not found", err.Error())
	case errors.Is(err, models.ErrProductAlreadyExists):
		response.Conflict(c, "product already exists", err.Error())
	default:
		response.InternalServerError(c, "internal server error", err.Error())
	}
}

func (h *ProductHandler) UploadProductImage(c *gin.Context) {
	productID := c.Param("id")

	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "invalid image upload", "field 'image' is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.BadRequest(c, "invalid image upload", "failed to open uploaded file")
		return
	}
	defer file.Close()

	product, err := h.productService.UploadImage(c.Request.Context(), services.UploadProductImageInput{
		ProductID:   productID,
		FileName:    fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
		File:        file,
	})
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "product image uploaded successfully", product)
}

func (h *ProductHandler) UploadProductGalleryImage(c *gin.Context) {
	productID := c.Param("id")

	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "invalid image upload", "field 'image' is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.BadRequest(c, "invalid image upload", "failed to open uploaded file")
		return
	}
	defer file.Close()

	image, err := h.productService.UploadGalleryImage(c.Request.Context(), services.UploadProductImageInput{
		ProductID:   productID,
		FileName:    fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
		File:        file,
	})
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "product image uploaded successfully", image)
}

func (h *ProductHandler) DeleteProductImage(c *gin.Context) {
	productID := c.Param("id")
	imageID := c.Param("image_id")

	if err := h.productService.DeleteImage(c.Request.Context(), productID, imageID); err != nil {
		h.handleProductError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) ReorderProductImages(c *gin.Context) {
	productID := c.Param("id")

	var req reorderProductImagesRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	images := make([]services.ReorderProductImageInput, 0, len(req.Images))
	for _, image := range req.Images {
		images = append(images, services.ReorderProductImageInput{
			ID:        image.ID,
			SortOrder: image.SortOrder,
		})
	}

	result, err := h.productService.ReorderImages(c.Request.Context(), services.ReorderProductImagesInput{
		ProductID: productID,
		Images:    images,
	})
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "product images reordered successfully", result)
}

func (h *ProductHandler) SetPrimaryProductImage(c *gin.Context) {
	productID := c.Param("id")
	imageID := c.Param("image_id")

	image, err := h.productService.SetPrimaryImage(c.Request.Context(), productID, imageID)
	if err != nil {
		h.handleProductError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "primary product image updated successfully", image)
}
