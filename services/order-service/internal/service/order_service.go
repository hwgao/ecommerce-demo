package service

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "ecommerce/services/order-service/internal/domain"
    "ecommerce/shared/pkg/events"
)

type OrderServiceImpl struct {
    orderRepo domain.OrderRepository
    eventBus  events.EventBus
}

func NewOrderService(orderRepo domain.OrderRepository, eventBus events.EventBus) *OrderServiceImpl {
    return &OrderServiceImpl{
        orderRepo: orderRepo,
        eventBus:  eventBus,
    }
}

func (s *OrderServiceImpl) CreateOrder(userID uuid.UUID, items []domain.OrderItem) (*domain.Order, error) {
    if len(items) == 0 {
        return nil, errors.New("order must contain at least one item")
    }

    // Calculate total amount
    var totalAmount float64
    for _, item := range items {
        totalAmount += item.Price * float64(item.Quantity)
    }

    order := &domain.Order{
        ID:          uuid.New(),
        UserID:      userID,
        Status:      "pending",
        TotalAmount: totalAmount,
        Currency:    "USD", // Default currency
        Items:       items,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := s.orderRepo.Create(order); err != nil {
        return nil, err
    }

    // Publish order created event
    event := events.OrderCreatedEvent{
        OrderID:     order.ID,
        UserID:      order.UserID,
        TotalAmount: order.TotalAmount,
        Currency:    order.Currency,
        Items:       items,
        Timestamp:   time.Now(),
    }
    s.eventBus.Publish("order.created", event)

    return order, nil
}

func (s *OrderServiceImpl) GetOrder(id uuid.UUID) (*domain.Order, error) {
    return s.orderRepo.GetByID(id)
}

func (s *OrderServiceImpl) GetUserOrders(userID uuid.UUID, limit, offset int) ([]*domain.Order, error) {
    return s.orderRepo.GetByUserID(userID, limit, offset)
}

func (s *OrderServiceImpl) UpdateOrderStatus(id uuid.UUID, status string) error {
    if err := s.orderRepo.UpdateStatus(id, status); err != nil {
        return err
    }

    // Publish order status updated event
    event := events.OrderStatusUpdatedEvent{
        OrderID:   id,
        Status:    status,
        Timestamp: time.Now(),
    }
    s.eventBus.Publish("order.status_updated", event)

    return nil
}

func (s *OrderServiceImpl) CancelOrder(id uuid.UUID) error {
    order, err := s.orderRepo.GetByID(id)
    if err != nil {
        return err
    }

    if order.Status == "shipped" || order.Status == "delivered" {
        return errors.New("cannot cancel shipped or delivered order")
    }

    return s.UpdateOrderStatus(id, "cancelled")
}
