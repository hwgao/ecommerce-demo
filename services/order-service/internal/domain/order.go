package domain

import (
    "time"
    "github.com/google/uuid"
)

type Order struct {
    ID         uuid.UUID   `json:"id" db:"id"`
    UserID     uuid.UUID   `json:"user_id" db:"user_id"`
    Status     string      `json:"status" db:"status"`
    TotalAmount float64    `json:"total_amount" db:"total_amount"`
    Currency   string      `json:"currency" db:"currency"`
    Items      []OrderItem `json:"items" db:"-"`
    CreatedAt  time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
    ID        uuid.UUID `json:"id" db:"id"`
    OrderID   uuid.UUID `json:"order_id" db:"order_id"`
    ProductID string    `json:"product_id" db:"product_id"`
    Quantity  int       `json:"quantity" db:"quantity"`
    Price     float64   `json:"price" db:"price"`
    Currency  string    `json:"currency" db:"currency"`
}

type OrderRepository interface {
    Create(order *Order) error
    GetByID(id uuid.UUID) (*Order, error)
    GetByUserID(userID uuid.UUID, limit, offset int) ([]*Order, error)
    Update(order *Order) error
    UpdateStatus(id uuid.UUID, status string) error
}

type OrderService interface {
    CreateOrder(userID uuid.UUID, items []OrderItem) (*Order, error)
    GetOrder(id uuid.UUID) (*Order, error)
    GetUserOrders(userID uuid.UUID, limit, offset int) ([]*Order, error)
    UpdateOrderStatus(id uuid.UUID, status string) error
    CancelOrder(id uuid.UUID) error
}
