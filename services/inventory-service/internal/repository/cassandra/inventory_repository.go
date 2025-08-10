package cassandra

import (
    "time"
    "github.com/gocql/gocql"
    "ecommerce/services/inventory-service/internal/domain"
)

type InventoryRepository struct {
    session *gocql.Session
}

func NewInventoryRepository(session *gocql.Session) *InventoryRepository {
    return &InventoryRepository{session: session}
}

func (r *InventoryRepository) CreateItem(item *domain.InventoryItem) error {
    item.Available = item.Quantity - item.Reserved
    item.LastUpdated = time.Now()
    
    query := `INSERT INTO inventory (product_id, quantity, reserved, available, last_updated) 
              VALUES (?, ?, ?, ?, ?)`
    
    return r.session.Query(query, item.ProductID, item.Quantity, 
        item.Reserved, item.Available, item.LastUpdated).Exec()
}

func (r *InventoryRepository) GetItem(productID string) (*domain.InventoryItem, error) {
    var item domain.InventoryItem
    
    query := `SELECT product_id, quantity, reserved, available, last_updated 
              FROM inventory WHERE product_id = ?`
    
    err := r.session.Query(query, productID).Scan(
        &item.ProductID, &item.Quantity, &item.Reserved, 
        &item.Available, &item.LastUpdated)
    
    if err != nil {
        return nil, err
    }
    
    return &item, nil
}

func (r *InventoryRepository) UpdateQuantity(productID string, quantity int) error {
    now := time.Now()
    
    query := `UPDATE inventory 
              SET quantity = ?, available = quantity - reserved, last_updated = ? 
              WHERE product_id = ?`
    
    return r.session.Query(query, quantity, now, productID).Exec()
}

func (r *InventoryRepository) ReserveQuantity(productID string, quantity int) error {
    now := time.Now()
    
    // Use lightweight transaction to ensure atomicity
    query := `UPDATE inventory 
              SET reserved = reserved + ?, available = available - ?, last_updated = ? 
              WHERE product_id = ? 
              IF available >= ?`
    
    applied, err := r.session.Query(query, quantity, quantity, now, productID, quantity).ScanCAS()
    if err != nil {
        return err
    }
    
    if !applied {
        return domain.ErrInsufficientStock
    }
    
    return nil
}

func (r *InventoryRepository) ReleaseReservation(productID string, quantity int) error {
    now := time.Now()
    
    query := `UPDATE inventory 
              SET reserved = reserved - ?, available = available + ?, last_updated = ? 
              WHERE product_id = ?`
    
    return r.session.Query(query, quantity, quantity, now, productID).Exec()
}
