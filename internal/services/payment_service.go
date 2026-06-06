package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type PaymentService interface {
	PayOrder(ctx context.Context, input PayOrderInput) (*models.Payment, error)
}

type PayOrderInput struct {
	UserID  string
	OrderID string
	Method  string
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo}
}

func (s *paymentService) PayOrder(ctx context.Context, input PayOrderInput) (*models.Payment, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidPaymentInput)
	}

	orderID := strings.TrimSpace(input.OrderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order id is required", models.ErrInvalidPaymentInput)
	}

	method := normalizePaymentMethod(input.Method)
	if !isValidPaymentMethod(method) {
		return nil, fmt.Errorf("%w: unsupported payment method", models.ErrInvalidPaymentInput)
	}

	transactionID, err := generatePaymentTransactionID()
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentRepo.CreateForOrder(ctx, repository.CreatePaymentInput{
		UserID:        userID,
		OrderID:       orderID,
		Provider:      models.PaymentProviderMock,
		PaymentMethod: method,
		TransactionID: transactionID,
		Status:        models.PaymentStatusPaid,
	})
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func normalizePaymentMethod(method string) string {
	return strings.ToLower(strings.TrimSpace(method))
}

func isValidPaymentMethod(method string) bool {
	switch method {
	case models.PaymentMethodBankTransfer,
		models.PaymentMethodCreditCard,
		models.PaymentMethodEWallet:
		return true
	default:
		return false
	}
}

func generatePaymentTransactionID() (string, error) {
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate transaction id: %w", err)
	}

	return fmt.Sprintf(
		"PAY-%s-%s",
		time.Now().UTC().Format("20060102150405"),
		strings.ToUpper(hex.EncodeToString(randomBytes)),
	), nil
}
