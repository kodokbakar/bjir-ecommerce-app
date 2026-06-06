package models

import "time"

const (
	PaymentProviderMock = "mock"

	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodCreditCard   = "credit_card"
	PaymentMethodEWallet      = "ewallet"

	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusFailed   = "failed"
	PaymentStatusExpired  = "expired"
	PaymentStatusRefunded = "refunded"
)

type Payment struct {
	ID            string     `json:"id"`
	OrderID       string     `json:"order_id"`
	Provider      string     `json:"provider"`
	PaymentMethod string     `json:"payment_method"`
	TransactionID string     `json:"transaction_id"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
