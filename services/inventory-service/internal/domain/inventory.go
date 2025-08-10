package domain

import (
    "time"
    "github.com/gocql/gocql"
)

type InventoryItem struct {
    ProductID   string    `json:"product_id" cql:"product_id"`
    Quantity    int       `json:"quantity" cql:"quantity"`
    Reserved    int       `json:"reserved" cql:"reserved"`
    Available   int       `json:"available" cql:"available"`
    LastUpdated time.Time `json:"last_updated" cql:"last_updated"`
}

type InventoryRepository interface {
    CreateItem(item *InventoryItem) error
    GetItem(productID string) (*InventoryItem, error)
    UpdateQuantity(productID string, quantity int) error
    ReserveQuantity(productID string, quantity int) error
    ReleaseReservation(productID string, quantity int) error
}

type InventoryService interface {
    AddStock(productID string, quantity int) error
    CheckAvailability(productID string) (*InventoryItem, error)
    ReserveStock(productID string, quantity int) error
    ReleaseReservation(productID string, quantity int) error
    UpdateStock(productID string, quantity int) error
}
