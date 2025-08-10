package service

import (
    "time"
    "github.com/google/uuid"
    "ecommerce/services/payment-service/internal/domain"
    "ecommerce/shared/pkg/events"
)

type PaymentServiceImpl struct {
    paymentRepo     domain.PaymentRepository
    paymentProvider domain.PaymentProvider
    eventBus        events.EventBus
}

func NewPaymentService(paymentRepo domain.PaymentRepository, 
                      paymentProvider domain.PaymentProvider,
                      eventBus events.EventBus) *PaymentServiceImpl {
    return &PaymentServiceImpl{
        paymentRepo:     paymentRepo,
        paymentProvider: paymentProvider,
        eventBus:        eventBus,
    }
}

func (s *PaymentServiceImpl) ProcessPayment(orderID, userID uuid.UUID, 
                                           amount float64, currency, paymentMethod string) (*domain.Payment, error) {
    payment := &domain.Payment{
        ID:            uuid.New(),
        OrderID:       orderID,
        UserID:        userID,
        Amount:        amount,
        Currency:      currency,
        PaymentMethod: paymentMethod,
        Status:        "pending",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if err := s.paymentRepo.Create(payment); err != nil {
        return nil, err
    }

    // Process payment with provider
    result, err := s.paymentProvider.ProcessPayment(amount, currency, paymentMethod)
    if err != nil {
        s.paymentRepo.UpdateStatus(payment.ID, "failed", err.Error())
        return nil, err
    }

    payment.ProviderID = result.ProviderID
    payment.Status = result.Status
    s.paymentRepo.UpdateStatus(payment.ID, result.Status, result.Response)

    // Publish payment event
    event := events.PaymentProcessedEvent{
        PaymentID: payment.ID,
        OrderID:   payment.OrderID,
        UserID:    payment.UserID,
        Amount:    payment.Amount,
        Currency:  payment.Currency,
        Status:    payment.Status,
        Timestamp: time.Now(),
    }
    s.eventBus.Publish("payment.processed", event)

    return payment, nil
}

func (s *PaymentServiceImpl) GetPayment(id uuid.UUID) (*domain.Payment, error) {
    return s.paymentRepo.GetByID(id)
}

func (s *PaymentServiceImpl) RefundPayment(paymentID uuid.UUID) error {
    payment, err := s.paymentRepo.GetByID(paymentID)
    if err != nil {
        return err
    }

    if payment.Status != "completed" {
        return errors.New("can only refund completed payments")
    }

    if err := s.paymentProvider.RefundPayment(payment.ProviderID, payment.Amount); err != nil {
        return err
    }

    s.paymentRepo.UpdateStatus(paymentID, "refunded", "Refund processed")
    
    // Publish refund event
    event := events.PaymentRefundedEvent{
        PaymentID: paymentID,
        OrderID:   payment.OrderID,
        Amount:    payment.Amount,
        Timestamp: time.Now(),
    }
    s.eventBus.Publish("payment.refunded", event)

    return nil
}
