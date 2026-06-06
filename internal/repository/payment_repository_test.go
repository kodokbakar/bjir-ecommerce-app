package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

func newPaymentRows(now time.Time) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"order_id",
		"provider",
		"payment_method",
		"transaction_id",
		"amount",
		"status",
		"paid_at",
		"created_at",
		"updated_at",
	}).AddRow(
		"payment-id",
		"order-id",
		models.PaymentProviderMock,
		models.PaymentMethodBankTransfer,
		"PAY-TEST",
		30000000.0,
		models.PaymentStatusPaid,
		now,
		now,
		now,
	)
}

func newPayableOrderRows(now time.Time) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"user_id",
		"order_number",
		"status",
		"total_amount",
		"shipping_address",
		"notes",
		"created_at",
		"updated_at",
	}).AddRow(
		"order-id",
		"user-id",
		"ORD-TEST",
		models.OrderStatusPending,
		30000000.0,
		"",
		"",
		now,
		now,
	)
}

func TestPaymentRepository_CreateForOrder_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepository(mock)

	now := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery("FROM orders").
		WithArgs("order-id", "user-id").
		WillReturnRows(newPayableOrderRows(now))

	mock.ExpectExec("UPDATE orders").
		WithArgs("order-id", models.OrderStatusPaid, models.OrderStatusPending).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectQuery("INSERT INTO payments").
		WithArgs(
			"order-id",
			models.PaymentProviderMock,
			models.PaymentMethodBankTransfer,
			"PAY-TEST",
			30000000.0,
			models.PaymentStatusPaid,
		).
		WillReturnRows(newPaymentRows(now))

	mock.ExpectCommit()

	payment, err := repo.CreateForOrder(context.Background(), CreatePaymentInput{
		UserID:        "user-id",
		OrderID:       "order-id",
		Provider:      models.PaymentProviderMock,
		PaymentMethod: models.PaymentMethodBankTransfer,
		TransactionID: "PAY-TEST",
		Status:        models.PaymentStatusPaid,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if payment.ID != "payment-id" {
		t.Fatalf("expected payment-id, got %s", payment.ID)
	}

	if payment.Status != models.PaymentStatusPaid {
		t.Fatalf("expected paid, got %s", payment.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPaymentRepository_CreateForOrder_OrderNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepository(mock)

	mock.ExpectBegin()

	mock.ExpectQuery("FROM orders").
		WithArgs("missing-order-id", "user-id").
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectRollback()

	payment, err := repo.CreateForOrder(context.Background(), CreatePaymentInput{
		UserID:        "user-id",
		OrderID:       "missing-order-id",
		Provider:      models.PaymentProviderMock,
		PaymentMethod: models.PaymentMethodEWallet,
		TransactionID: "PAY-TEST",
		Status:        models.PaymentStatusPaid,
	})
	if payment != nil {
		t.Fatal("expected nil payment")
	}

	if !errors.Is(err, models.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPaymentRepository_CreateForOrder_OrderNotPayable(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepository(mock)

	now := time.Now()

	paidOrderRows := pgxmock.NewRows([]string{
		"id",
		"user_id",
		"order_number",
		"status",
		"total_amount",
		"shipping_address",
		"notes",
		"created_at",
		"updated_at",
	}).AddRow(
		"order-id",
		"user-id",
		"ORD-TEST",
		models.OrderStatusPaid,
		30000000.0,
		"",
		"",
		now,
		now,
	)

	mock.ExpectBegin()

	mock.ExpectQuery("FROM orders").
		WithArgs("order-id", "user-id").
		WillReturnRows(paidOrderRows)

	mock.ExpectRollback()

	payment, err := repo.CreateForOrder(context.Background(), CreatePaymentInput{
		UserID:        "user-id",
		OrderID:       "order-id",
		Provider:      models.PaymentProviderMock,
		PaymentMethod: models.PaymentMethodCreditCard,
		TransactionID: "PAY-TEST",
		Status:        models.PaymentStatusPaid,
	})
	if payment != nil {
		t.Fatal("expected nil payment")
	}

	if !errors.Is(err, models.ErrOrderNotPayable) {
		t.Fatalf("expected ErrOrderNotPayable, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
