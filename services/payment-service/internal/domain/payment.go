package domain

import (
    "time"
    "github.com/google/uuid"
)

type Payment struct {
    ID              uuid.UUID `json:"id" db:"id"`
    OrderID         uuid.UUID `json:"order_id" db:"order_id"`
    UserID          uuid.UUID `json:"user_id" db:"user_id"`
    Amount          float64   `json:"amount" db:"amount"`
    Currency        string    `json:"currency" db:"currency"`
    PaymentMethod   string    `json:"payment_method" db:"payment_method"`
    Status          string    `json:"status" db:"status"`
    ProviderID      string    `json:"provider_id" db:"provider_id"`
    ProviderResponse string   `json:"-" db:"provider_response"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type PaymentRepository interface {
    Create(payment *Payment) error
    GetByID(id uuid.UUID) (*Payment, error)
    GetByOrderID(orderID uuid.UUID) (*Payment, error)
    UpdateStatus(id uuid.UUID, status, providerResponse string) error
}

type PaymentProvider interface {
    ProcessPayment(amount float64, currency, paymentMethod string) (*PaymentResult, error)
    RefundPayment(providerID string, amount float64) error
}

type PaymentResult struct {
    ProviderID string
    Status     string
    Response   string
}

type PaymentService interface {
    ProcessPayment(orderID, userID uuid.UUID, amount float64, currency, paymentMethod string) (*Payment, error)
    GetPayment(id uuid.UUID) (*Payment, error)
    RefundPayment(paymentID uuid.UUID) error
}
