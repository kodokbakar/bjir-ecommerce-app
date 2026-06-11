package models

import "time"

type Product struct {
	ID          string         `json:"id"`
	CategoryID  string         `json:"category_id"`
	Category    *Category      `json:"category,omitempty"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	ImageURL    string         `json:"image_url"`
	Images      []ProductImage `json:"images,omitempty"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
