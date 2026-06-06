package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type fakePaymentRepository struct {
	createForOrderFunc func(ctx context.Context, input repository.CreatePaymentInput) (*models.Payment, error)
}

func newFakePaymentRepository() *fakePaymentRepository {
	return &fakePaymentRepository{
		createForOrderFunc: func(ctx context.Context, input repository.CreatePaymentInput) (*models.Payment, error) {
			return &models.Payment{
				ID:            "payment-id",
				OrderID:       input.OrderID,
				Provider:      input.Provider,
				PaymentMethod: input.PaymentMethod,
				TransactionID: input.TransactionID,
				Amount:        30000000,
				Status:        input.Status,
			}, nil
		},
	}
}

func (f *fakePaymentRepository) CreateForOrder(ctx context.Context, input repository.CreatePaymentInput) (*models.Payment, error) {
	return f.createForOrderFunc(ctx, input)
}

func TestPaymentService_PayOrder_Success(t *testing.T) {
	repo := newFakePaymentRepository()

	var receivedInput repository.CreatePaymentInput
	repo.createForOrderFunc = func(ctx context.Context, input repository.CreatePaymentInput) (*models.Payment, error) {
		receivedInput = input

		return &models.Payment{
			ID:            "payment-id",
			OrderID:       input.OrderID,
			Provider:      input.Provider,
			PaymentMethod: input.PaymentMethod,
			TransactionID: input.TransactionID,
			Amount:        30000000,
			Status:        input.Status,
		}, nil
	}

	service := NewPaymentService(repo)

	payment, err := service.PayOrder(context.Background(), PayOrderInput{
		UserID:  "user-id",
		OrderID: "order-id",
		Method:  "bank_transfer",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if payment.ID != "payment-id" {
		t.Fatalf("expected payment-id, got %s", payment.ID)
	}

	if receivedInput.UserID != "user-id" {
		t.Fatalf("expected user-id, got %s", receivedInput.UserID)
	}

	if receivedInput.OrderID != "order-id" {
		t.Fatalf("expected order-id, got %s", receivedInput.OrderID)
	}

	if receivedInput.Provider != models.PaymentProviderMock {
		t.Fatalf("expected provider mock, got %s", receivedInput.Provider)
	}

	if receivedInput.PaymentMethod != models.PaymentMethodBankTransfer {
		t.Fatalf("expected bank_transfer, got %s", receivedInput.PaymentMethod)
	}

	if receivedInput.Status != models.PaymentStatusPaid {
		t.Fatalf("expected paid, got %s", receivedInput.Status)
	}

	if receivedInput.TransactionID == "" {
		t.Fatal("expected transaction id to be generated")
	}

	if !strings.HasPrefix(receivedInput.TransactionID, "PAY-") {
		t.Fatalf("expected transaction id prefix PAY-, got %s", receivedInput.TransactionID)
	}
}

func TestPaymentService_PayOrder_InvalidInput(t *testing.T) {
	service := NewPaymentService(newFakePaymentRepository())

	tests := []struct {
		name  string
		input PayOrderInput
	}{
		{
			name: "empty user id",
			input: PayOrderInput{
				UserID:  "",
				OrderID: "order-id",
				Method:  "bank_transfer",
			},
		},
		{
			name: "empty order id",
			input: PayOrderInput{
				UserID:  "user-id",
				OrderID: "",
				Method:  "bank_transfer",
			},
		},
		{
			name: "invalid method",
			input: PayOrderInput{
				UserID:  "user-id",
				OrderID: "order-id",
				Method:  "cash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.PayOrder(context.Background(), tt.input)
			if !errors.Is(err, models.ErrInvalidPaymentInput) {
				t.Fatalf("expected ErrInvalidPaymentInput, got %v", err)
			}
		})
	}
}

func TestPaymentService_PayOrder_OrderNotFound(t *testing.T) {
	repo := newFakePaymentRepository()
	repo.createForOrderFunc = func(ctx context.Context, input repository.CreatePaymentInput) (*models.Payment, error) {
		return nil, models.ErrOrderNotFound
	}

	service := NewPaymentService(repo)

	_, err := service.PayOrder(context.Background(), PayOrderInput{
		UserID:  "user-id",
		OrderID: "missing-order-id",
		Method:  "ewallet",
	})

	if !errors.Is(err, models.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestPaymentService_PayOrder_OrderNotPayable(t *testing.T) {
	repo := newFakePaymentRepository()
	repo.createForOrderFunc = func(ctx context.Context, input repository.CreatePaymentInput) (*models.Payment, error) {
		return nil, models.ErrOrderNotPayable
	}

	service := NewPaymentService(repo)

	_, err := service.PayOrder(context.Background(), PayOrderInput{
		UserID:  "user-id",
		OrderID: "order-id",
		Method:  "credit_card",
	})

	if !errors.Is(err, models.ErrOrderNotPayable) {
		t.Fatalf("expected ErrOrderNotPayable, got %v", err)
	}
}
