package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/database"
)

const defaultSeedPassword = "password123"

type seedUser struct {
	Name  string
	Email string
	Role  string
}

type seedCategory struct {
	Name        string
	Slug        string
	Description string
	ImageURL    string
	ParentSlug  string
}

type seedProduct struct {
	CategorySlug string
	Name         string
	Slug         string
	Description  string
	Price        int64
	Stock        int
	ImageURL     string
}

type seededProduct struct {
	ID    string
	Name  string
	Price int64
}

type seedCartItem struct {
	UserEmail   string
	ProductSlug string
	Quantity    int
}

type seedOrderItem struct {
	ProductSlug string
	Quantity    int
}

type seedPayment struct {
	Provider      string
	PaymentMethod string
	TransactionID string
	Status        string
	PaidAt        *time.Time
}

type seedOrder struct {
	OrderNumber     string
	UserEmail       string
	Status          string
	ShippingAddress string
	Notes           string
	Items           []seedOrderItem
	Payment         *seedPayment
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("begin transaction: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(context.Background())
		}
	}()

	if err := seed(ctx, tx); err != nil {
		log.Fatalf("seed database: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("commit transaction: %v", err)
	}
	committed = true

	log.Println("database seeded successfully")
	log.Println("admin login: admin@bjir.com / password123")
	log.Println("customer login: customer1@bjir.com / password123")
	log.Println("customer login: customer2@bjir.com / password123")
}

func seed(ctx context.Context, tx pgx.Tx) error {
	userIDs, err := seedUsers(ctx, tx)
	if err != nil {
		return err
	}

	categoryIDs, err := seedCategories(ctx, tx)
	if err != nil {
		return err
	}

	products, err := seedProducts(ctx, tx, categoryIDs)
	if err != nil {
		return err
	}

	if err := seedProductImages(ctx, tx, products); err != nil {
		return err
	}

	if err := seedCarts(ctx, tx, userIDs, products); err != nil {
		return err
	}

	if err := seedOrders(ctx, tx, userIDs, products); err != nil {
		return err
	}

	return nil
}

func seedUsers(ctx context.Context, tx pgx.Tx) (map[string]string, error) {
	seedPassword := getSeedPassword()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash seed password: %w", err)
	}

	users := []seedUser{
		{Name: "Admin Bjir", Email: "admin@bjir.com", Role: "admin"},
		{Name: "Alya Customer", Email: "customer1@bjir.com", Role: "customer"},
		{Name: "Bima Customer", Email: "customer2@bjir.com", Role: "customer"},
	}

	userIDs := make(map[string]string, len(users))

	for _, user := range users {
		var id string
		err := tx.QueryRow(ctx, `
			INSERT INTO users (
				name,
				email,
				password_hash,
				role,
				is_active
			)
			VALUES ($1, $2, $3, $4, TRUE)
			ON CONFLICT (email) DO UPDATE
			SET name = EXCLUDED.name,
				password_hash = EXCLUDED.password_hash,
				role = EXCLUDED.role,
				is_active = TRUE,
				updated_at = NOW()
			RETURNING id::text
		`, user.Name, user.Email, string(passwordHash), user.Role).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("seed user %s: %w", user.Email, err)
		}

		userIDs[user.Email] = id
	}

	return userIDs, nil
}

func seedCategories(ctx context.Context, tx pgx.Tx) (map[string]string, error) {
	categories := []seedCategory{
		{
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Gadgets, devices, and daily tech essentials.",
			ImageURL:    "https://placehold.co/800x600?text=Electronics",
		},
		{
			Name:        "Laptops",
			Slug:        "laptops",
			Description: "Portable computers for work, school, and gaming.",
			ImageURL:    "https://placehold.co/800x600?text=Laptops",
			ParentSlug:  "electronics",
		},
		{
			Name:        "Smartphones",
			Slug:        "smartphones",
			Description: "Android and iOS smartphones with modern features.",
			ImageURL:    "https://placehold.co/800x600?text=Smartphones",
			ParentSlug:  "electronics",
		},
		{
			Name:        "Fashion",
			Slug:        "fashion",
			Description: "Everyday fashion and accessories.",
			ImageURL:    "https://placehold.co/800x600?text=Fashion",
		},
		{
			Name:        "Men's Fashion",
			Slug:        "mens-fashion",
			Description: "Clothing and accessories for men.",
			ImageURL:    "https://placehold.co/800x600?text=Mens+Fashion",
			ParentSlug:  "fashion",
		},
		{
			Name:        "Women's Fashion",
			Slug:        "womens-fashion",
			Description: "Clothing and accessories for women.",
			ImageURL:    "https://placehold.co/800x600?text=Womens+Fashion",
			ParentSlug:  "fashion",
		},
		{
			Name:        "Home Living",
			Slug:        "home-living",
			Description: "Furniture, lighting, and home improvement.",
			ImageURL:    "https://placehold.co/800x600?text=Home+Living",
		},
		{
			Name:        "Kitchen",
			Slug:        "kitchen",
			Description: "Kitchen tools and cooking appliances.",
			ImageURL:    "https://placehold.co/800x600?text=Kitchen",
			ParentSlug:  "home-living",
		},
	}

	categoryIDs := make(map[string]string, len(categories))

	// Pass 1: upsert all categories without parent_id first.
	// This makes seed data independent from category slice order.
	for _, category := range categories {
		var id string
		err := tx.QueryRow(ctx, `
			INSERT INTO categories (
				name,
				slug,
				description,
				image_url,
				parent_id,
				is_active
			)
			VALUES ($1, $2, $3, $4, NULL, TRUE)
			ON CONFLICT (slug) DO UPDATE
			SET name = EXCLUDED.name,
				description = EXCLUDED.description,
				image_url = EXCLUDED.image_url,
				parent_id = NULL,
				is_active = TRUE,
				updated_at = NOW()
			RETURNING id::text
		`, category.Name, category.Slug, category.Description, category.ImageURL).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("seed category %s: %w", category.Slug, err)
		}

		categoryIDs[category.Slug] = id
	}

	// Pass 2: resolve parent_id after every category slug already has an ID.
	for _, category := range categories {
		var parentID *string

		if category.ParentSlug != "" {
			id, ok := categoryIDs[category.ParentSlug]
			if !ok {
				return nil, fmt.Errorf("parent category %s not found for %s", category.ParentSlug, category.Slug)
			}

			parentID = &id
		}

		if _, err := tx.Exec(ctx, `
			UPDATE categories
			SET parent_id = $2,
				updated_at = NOW()
			WHERE slug = $1
		`, category.Slug, parentID); err != nil {
			return nil, fmt.Errorf("set parent category %s: %w", category.Slug, err)
		}
	}

	return categoryIDs, nil
}

func seedProducts(ctx context.Context, tx pgx.Tx, categoryIDs map[string]string) (map[string]seededProduct, error) {
	products := []seedProduct{
		{
			CategorySlug: "laptops",
			Name:         "Acer Swift Go 14",
			Slug:         "acer-swift-go-14",
			Description:  "Lightweight 14-inch laptop for students and remote workers.",
			Price:        4999000,
			Stock:        15,
			ImageURL:     "https://placehold.co/900x700?text=Acer+Swift+Go+14",
		},
		{
			CategorySlug: "laptops",
			Name:         "Lenovo ThinkPad X280",
			Slug:         "lenovo-thinkpad-x280",
			Description:  "Compact business laptop with durable keyboard and solid battery life.",
			Price:        3850000,
			Stock:        9,
			ImageURL:     "https://placehold.co/900x700?text=ThinkPad+X280",
		},
		{
			CategorySlug: "laptops",
			Name:         "ASUS Vivobook 15",
			Slug:         "asus-vivobook-15",
			Description:  "Daily productivity laptop with large display and fast SSD.",
			Price:        5299000,
			Stock:        11,
			ImageURL:     "https://placehold.co/900x700?text=ASUS+Vivobook+15",
		},
		{
			CategorySlug: "smartphones",
			Name:         "Samsung Galaxy A55",
			Slug:         "samsung-galaxy-a55",
			Description:  "Mid-range Android phone with AMOLED display and long battery life.",
			Price:        4399000,
			Stock:        24,
			ImageURL:     "https://placehold.co/900x700?text=Galaxy+A55",
		},
		{
			CategorySlug: "smartphones",
			Name:         "iPhone SE 2022",
			Slug:         "iphone-se-2022",
			Description:  "Compact iPhone with reliable performance and iOS support.",
			Price:        4999000,
			Stock:        8,
			ImageURL:     "https://placehold.co/900x700?text=iPhone+SE+2022",
		},
		{
			CategorySlug: "smartphones",
			Name:         "Redmi Note 13",
			Slug:         "redmi-note-13",
			Description:  "Affordable phone with large display, good camera, and big battery.",
			Price:        2499000,
			Stock:        31,
			ImageURL:     "https://placehold.co/900x700?text=Redmi+Note+13",
		},
		{
			CategorySlug: "mens-fashion",
			Name:         "Oxford Cotton Shirt",
			Slug:         "oxford-cotton-shirt",
			Description:  "Classic cotton shirt for office and casual outfits.",
			Price:        189000,
			Stock:        60,
			ImageURL:     "https://placehold.co/900x700?text=Oxford+Shirt",
		},
		{
			CategorySlug: "mens-fashion",
			Name:         "Slim Fit Chino Pants",
			Slug:         "slim-fit-chino-pants",
			Description:  "Comfortable chino pants with clean slim silhouette.",
			Price:        249000,
			Stock:        48,
			ImageURL:     "https://placehold.co/900x700?text=Chino+Pants",
		},
		{
			CategorySlug: "mens-fashion",
			Name:         "Canvas Sneakers",
			Slug:         "canvas-sneakers",
			Description:  "Everyday sneakers with lightweight canvas upper.",
			Price:        329000,
			Stock:        36,
			ImageURL:     "https://placehold.co/900x700?text=Canvas+Sneakers",
		},
		{
			CategorySlug: "womens-fashion",
			Name:         "Linen Midi Dress",
			Slug:         "linen-midi-dress",
			Description:  "Breathable midi dress for warm weather and casual occasions.",
			Price:        279000,
			Stock:        35,
			ImageURL:     "https://placehold.co/900x700?text=Linen+Midi+Dress",
		},
		{
			CategorySlug: "womens-fashion",
			Name:         "Everyday Tote Bag",
			Slug:         "everyday-tote-bag",
			Description:  "Roomy tote bag for laptop, books, and daily essentials.",
			Price:        159000,
			Stock:        44,
			ImageURL:     "https://placehold.co/900x700?text=Tote+Bag",
		},
		{
			CategorySlug: "womens-fashion",
			Name:         "Soft Knit Cardigan",
			Slug:         "soft-knit-cardigan",
			Description:  "Soft cardigan with relaxed fit for layering.",
			Price:        219000,
			Stock:        28,
			ImageURL:     "https://placehold.co/900x700?text=Knit+Cardigan",
		},
		{
			CategorySlug: "home-living",
			Name:         "Minimalist Desk Lamp",
			Slug:         "minimalist-desk-lamp",
			Description:  "Adjustable LED desk lamp with warm and cool lighting modes.",
			Price:        175000,
			Stock:        50,
			ImageURL:     "https://placehold.co/900x700?text=Desk+Lamp",
		},
		{
			CategorySlug: "home-living",
			Name:         "Ergonomic Office Chair",
			Slug:         "ergonomic-office-chair",
			Description:  "Comfortable chair with lumbar support for long work sessions.",
			Price:        1499000,
			Stock:        14,
			ImageURL:     "https://placehold.co/900x700?text=Office+Chair",
		},
		{
			CategorySlug: "home-living",
			Name:         "Floating Wall Shelf",
			Slug:         "floating-wall-shelf",
			Description:  "Simple wall shelf for books, plants, and display items.",
			Price:        99000,
			Stock:        70,
			ImageURL:     "https://placehold.co/900x700?text=Wall+Shelf",
		},
		{
			CategorySlug: "kitchen",
			Name:         "Stainless Steel Pan",
			Slug:         "stainless-steel-pan",
			Description:  "Durable 28cm frying pan for daily cooking.",
			Price:        199000,
			Stock:        42,
			ImageURL:     "https://placehold.co/900x700?text=Steel+Pan",
		},
		{
			CategorySlug: "kitchen",
			Name:         "Electric Coffee Grinder",
			Slug:         "electric-coffee-grinder",
			Description:  "Compact grinder for fresh coffee at home.",
			Price:        449000,
			Stock:        22,
			ImageURL:     "https://placehold.co/900x700?text=Coffee+Grinder",
		},
		{
			CategorySlug: "kitchen",
			Name:         "Ceramic Dinner Plate Set",
			Slug:         "ceramic-dinner-plate-set",
			Description:  "Six-piece ceramic plate set for family dining.",
			Price:        129000,
			Stock:        55,
			ImageURL:     "https://placehold.co/900x700?text=Plate+Set",
		},
	}

	seeded := make(map[string]seededProduct, len(products))

	for _, product := range products {
		categoryID, ok := categoryIDs[product.CategorySlug]
		if !ok {
			return nil, fmt.Errorf("category %s not found for product %s", product.CategorySlug, product.Slug)
		}

		var id string
		err := tx.QueryRow(ctx, `
			INSERT INTO products (
				category_id,
				name,
				slug,
				description,
				price,
				stock,
				image_url,
				is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)
			ON CONFLICT (slug) DO UPDATE
			SET category_id = EXCLUDED.category_id,
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				price = EXCLUDED.price,
				stock = EXCLUDED.stock,
				image_url = EXCLUDED.image_url,
				is_active = TRUE,
				updated_at = NOW()
			RETURNING id::text
		`, categoryID, product.Name, product.Slug, product.Description, product.Price, product.Stock, product.ImageURL).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("seed product %s: %w", product.Slug, err)
		}

		seeded[product.Slug] = seededProduct{
			ID:    id,
			Name:  product.Name,
			Price: product.Price,
		}
	}

	return seeded, nil
}

func seedProductImages(ctx context.Context, tx pgx.Tx, products map[string]seededProduct) error {
	for slug, product := range products {
		imageURL := fmt.Sprintf("https://placehold.co/900x700?text=%s", slug)

		var existingID string
		err := tx.QueryRow(ctx, `
			SELECT id::text
			FROM product_images
			WHERE product_id = $1
			AND is_primary = TRUE
			ORDER BY created_at ASC
			LIMIT 1
		`, product.ID).Scan(&existingID)

		switch {
		case err == nil:
			if _, err := tx.Exec(ctx, `
				UPDATE product_images
				SET image_url = $2,
					sort_order = 0,
					updated_at = NOW()
				WHERE id = $1
			`, existingID, imageURL); err != nil {
				return fmt.Errorf("update product image %s: %w", slug, err)
			}
		case errors.Is(err, pgx.ErrNoRows):
			if _, err := tx.Exec(ctx, `
				INSERT INTO product_images (
					product_id,
					image_url,
					sort_order,
					is_primary
				)
				VALUES ($1, $2, 0, TRUE)
			`, product.ID, imageURL); err != nil {
				return fmt.Errorf("insert product image %s: %w", slug, err)
			}
		default:
			return fmt.Errorf("find product image %s: %w", slug, err)
		}

		if _, err := tx.Exec(ctx, `
			UPDATE products
			SET image_url = $2,
				updated_at = NOW()
			WHERE id = $1
		`, product.ID, imageURL); err != nil {
			return fmt.Errorf("sync product image_url %s: %w", slug, err)
		}
	}

	return nil
}

func seedCarts(ctx context.Context, tx pgx.Tx, userIDs map[string]string, products map[string]seededProduct) error {
	items := []seedCartItem{
		{UserEmail: "customer1@bjir.com", ProductSlug: "redmi-note-13", Quantity: 1},
		{UserEmail: "customer1@bjir.com", ProductSlug: "oxford-cotton-shirt", Quantity: 2},
		{UserEmail: "customer1@bjir.com", ProductSlug: "minimalist-desk-lamp", Quantity: 1},
		{UserEmail: "customer2@bjir.com", ProductSlug: "acer-swift-go-14", Quantity: 1},
		{UserEmail: "customer2@bjir.com", ProductSlug: "everyday-tote-bag", Quantity: 1},
		{UserEmail: "customer2@bjir.com", ProductSlug: "electric-coffee-grinder", Quantity: 1},
		{UserEmail: "customer2@bjir.com", ProductSlug: "ceramic-dinner-plate-set", Quantity: 3},
	}

	for _, item := range items {
		userID, ok := userIDs[item.UserEmail]
		if !ok {
			return fmt.Errorf("user %s not found for cart", item.UserEmail)
		}

		product, ok := products[item.ProductSlug]
		if !ok {
			return fmt.Errorf("product %s not found for cart", item.ProductSlug)
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO carts (
				user_id,
				product_id,
				quantity
			)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, product_id) DO UPDATE
			SET quantity = EXCLUDED.quantity,
				updated_at = NOW()
		`, userID, product.ID, item.Quantity); err != nil {
			return fmt.Errorf("seed cart item %s/%s: %w", item.UserEmail, item.ProductSlug, err)
		}
	}

	return nil
}

func seedOrders(ctx context.Context, tx pgx.Tx, userIDs map[string]string, products map[string]seededProduct) error {
	orders := []seedOrder{
		{
			OrderNumber:     "SEED-ORDER-0001",
			UserEmail:       "customer1@bjir.com",
			Status:          "pending",
			ShippingAddress: "Jl. Melati No. 12, Malang, Jawa Timur",
			Notes:           "Seed pending order for local development.",
			Items: []seedOrderItem{
				{ProductSlug: "redmi-note-13", Quantity: 1},
				{ProductSlug: "canvas-sneakers", Quantity: 1},
			},
			Payment: &seedPayment{
				Provider:      "midtrans",
				PaymentMethod: "bank_transfer",
				TransactionID: "SEED-PAY-0001",
				Status:        "pending",
			},
		},
		{
			OrderNumber:     "SEED-ORDER-0002",
			UserEmail:       "customer1@bjir.com",
			Status:          "paid",
			ShippingAddress: "Jl. Ijen No. 22, Malang, Jawa Timur",
			Notes:           "Paid order waiting for fulfillment.",
			Items: []seedOrderItem{
				{ProductSlug: "oxford-cotton-shirt", Quantity: 2},
				{ProductSlug: "slim-fit-chino-pants", Quantity: 1},
				{ProductSlug: "everyday-tote-bag", Quantity: 1},
			},
			Payment: &seedPayment{
				Provider:      "midtrans",
				PaymentMethod: "qris",
				TransactionID: "SEED-PAY-0002",
				Status:        "paid",
				PaidAt:        daysAgo(5),
			},
		},
		{
			OrderNumber:     "SEED-ORDER-0003",
			UserEmail:       "customer2@bjir.com",
			Status:          "shipped",
			ShippingAddress: "Jl. Sudirman No. 8, Surabaya, Jawa Timur",
			Notes:           "Courier has picked up the package.",
			Items: []seedOrderItem{
				{ProductSlug: "acer-swift-go-14", Quantity: 1},
				{ProductSlug: "minimalist-desk-lamp", Quantity: 1},
			},
			Payment: &seedPayment{
				Provider:      "midtrans",
				PaymentMethod: "credit_card",
				TransactionID: "SEED-PAY-0003",
				Status:        "paid",
				PaidAt:        daysAgo(4),
			},
		},
		{
			OrderNumber:     "SEED-ORDER-0004",
			UserEmail:       "customer2@bjir.com",
			Status:          "delivered",
			ShippingAddress: "Jl. Asia Afrika No. 18, Bandung, Jawa Barat",
			Notes:           "Delivered seed order.",
			Items: []seedOrderItem{
				{ProductSlug: "electric-coffee-grinder", Quantity: 1},
				{ProductSlug: "stainless-steel-pan", Quantity: 1},
				{ProductSlug: "ceramic-dinner-plate-set", Quantity: 2},
			},
			Payment: &seedPayment{
				Provider:      "midtrans",
				PaymentMethod: "bank_transfer",
				TransactionID: "SEED-PAY-0004",
				Status:        "paid",
				PaidAt:        daysAgo(9),
			},
		},
		{
			OrderNumber:     "SEED-ORDER-0005",
			UserEmail:       "customer1@bjir.com",
			Status:          "cancelled",
			ShippingAddress: "Jl. Diponegoro No. 15, Yogyakarta",
			Notes:           "Cancelled seed order with failed payment.",
			Items: []seedOrderItem{
				{ProductSlug: "soft-knit-cardigan", Quantity: 1},
				{ProductSlug: "linen-midi-dress", Quantity: 1},
			},
			Payment: &seedPayment{
				Provider:      "midtrans",
				PaymentMethod: "virtual_account",
				TransactionID: "SEED-PAY-0005",
				Status:        "failed",
			},
		},
		{
			OrderNumber:     "SEED-ORDER-0006",
			UserEmail:       "customer2@bjir.com",
			Status:          "processing",
			ShippingAddress: "Jl. Gatot Subroto No. 20, Jakarta Selatan",
			Notes:           "Processing order for dashboard state coverage.",
			Items: []seedOrderItem{
				{ProductSlug: "samsung-galaxy-a55", Quantity: 1},
				{ProductSlug: "floating-wall-shelf", Quantity: 2},
			},
			Payment: &seedPayment{
				Provider:      "midtrans",
				PaymentMethod: "e_wallet",
				TransactionID: "SEED-PAY-0006",
				Status:        "paid",
				PaidAt:        daysAgo(2),
			},
		},
	}

	for _, order := range orders {
		if err := seedOrderWithItemsAndPayment(ctx, tx, order, userIDs, products); err != nil {
			return err
		}
	}

	return nil
}

func seedOrderWithItemsAndPayment(
	ctx context.Context,
	tx pgx.Tx,
	order seedOrder,
	userIDs map[string]string,
	products map[string]seededProduct,
) error {
	userID, ok := userIDs[order.UserEmail]
	if !ok {
		return fmt.Errorf("user %s not found for order %s", order.UserEmail, order.OrderNumber)
	}

	totalAmount, err := calculateOrderTotal(order, products)
	if err != nil {
		return err
	}

	var orderID string
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			user_id,
			order_number,
			status,
			total_amount,
			shipping_address,
			notes
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (order_number) DO UPDATE
		SET user_id = EXCLUDED.user_id,
			status = EXCLUDED.status,
			total_amount = EXCLUDED.total_amount,
			shipping_address = EXCLUDED.shipping_address,
			notes = EXCLUDED.notes,
			updated_at = NOW()
		RETURNING id::text
	`, userID, order.OrderNumber, order.Status, totalAmount, order.ShippingAddress, order.Notes).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("seed order %s: %w", order.OrderNumber, err)
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM order_items
		WHERE order_id = $1
	`, orderID); err != nil {
		return fmt.Errorf("clear order items %s: %w", order.OrderNumber, err)
	}

	for _, item := range order.Items {
		product, ok := products[item.ProductSlug]
		if !ok {
			return fmt.Errorf("product %s not found for order %s", item.ProductSlug, order.OrderNumber)
		}

		subtotal := product.Price * int64(item.Quantity)

		if _, err := tx.Exec(ctx, `
			INSERT INTO order_items (
				order_id,
				product_id,
				product_name,
				quantity,
				price,
				subtotal
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, orderID, product.ID, product.Name, item.Quantity, product.Price, subtotal); err != nil {
			return fmt.Errorf("seed order item %s/%s: %w", order.OrderNumber, item.ProductSlug, err)
		}
	}

	if order.Payment == nil {
		if _, err := tx.Exec(ctx, `
			DELETE FROM payments
			WHERE order_id = $1
		`, orderID); err != nil {
			return fmt.Errorf("delete payment %s: %w", order.OrderNumber, err)
		}

		return nil
	}

	if err := seedPaymentForOrder(ctx, tx, orderID, order.OrderNumber, totalAmount, *order.Payment); err != nil {
		return err
	}

	return nil
}

func seedPaymentForOrder(ctx context.Context, tx pgx.Tx, orderID string, orderNumber string, amount int64, payment seedPayment) error {
	var paidAt any
	if payment.PaidAt != nil {
		paidAt = *payment.PaidAt
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO payments (
			order_id,
			provider,
			payment_method,
			transaction_id,
			amount,
			status,
			paid_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (order_id) DO UPDATE
		SET provider = EXCLUDED.provider,
			payment_method = EXCLUDED.payment_method,
			transaction_id = EXCLUDED.transaction_id,
			amount = EXCLUDED.amount,
			status = EXCLUDED.status,
			paid_at = EXCLUDED.paid_at,
			updated_at = NOW()
	`, orderID, payment.Provider, payment.PaymentMethod, payment.TransactionID, amount, payment.Status, paidAt); err != nil {
		return fmt.Errorf("seed payment %s: %w", orderNumber, err)
	}

	return nil
}

func calculateOrderTotal(order seedOrder, products map[string]seededProduct) (int64, error) {
	var total int64

	for _, item := range order.Items {
		product, ok := products[item.ProductSlug]
		if !ok {
			return 0, fmt.Errorf("product %s not found for order %s", item.ProductSlug, order.OrderNumber)
		}

		total += product.Price * int64(item.Quantity)
	}

	return total, nil
}

func daysAgo(days int) *time.Time {
	value := time.Now().AddDate(0, 0, -days)
	return &value
}

func getSeedPassword() string {
	password := os.Getenv("SEED_PASSWORD")
	if password == "" {
		return defaultSeedPassword
	}

	return password
}
