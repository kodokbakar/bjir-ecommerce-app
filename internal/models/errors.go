package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryHasProducts   = errors.New("category has related products")
	ErrCategoryHasChildren   = errors.New("category has child categories")
	ErrInvalidCategoryInput  = errors.New("invalid category input")

	ErrProductNotFound      = errors.New("product not found")
	ErrProductImageNotFound = errors.New("product image not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidProductInput  = errors.New("invalid product input")

	ErrCartItemNotFound = errors.New("cart item not found")
	ErrInvalidCartInput = errors.New("invalid cart input")
	ErrCartEmpty        = errors.New("cart is empty")

	ErrInvalidOrderStatus           = errors.New("invalid order status")
	ErrInvalidOrderStatusTransition = errors.New("invalid order status transition")
	ErrOrderNotFound                = errors.New("order not found")
	ErrInvalidOrderInput            = errors.New("invalid order input")
	ErrInsufficientStock            = errors.New("insufficient product stock")

	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentAlreadyExists = errors.New("payment already exists")
	ErrInvalidPaymentInput  = errors.New("invalid payment input")
	ErrOrderNotPayable      = errors.New("order is not payable")
)
