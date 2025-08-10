package domain

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/google/uuid"
)

type Notification struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    UserID    uuid.UUID         `json:"user_id" bson:"user_id"`
    Type      string            `json:"type" bson:"type"` // email, sms, push
    Channel   string            `json:"channel" bson:"channel"`
    Subject   string            `json:"subject" bson:"subject"`
    Content   string            `json:"content" bson:"content"`
    Status    string            `json:"status" bson:"status"`
    Metadata  map[string]interface{} `json:"metadata" bson:"metadata"`
    CreatedAt time.Time         `json:"created_at" bson:"created_at"`
    SentAt    *time.Time        `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
}

type NotificationRepository interface {
    Create(notification *Notification) error
    GetByID(id string) (*Notification, error)
    GetByUserID(userID uuid.UUID, limit, offset int) ([]*Notification, error)
    UpdateStatus(id string, status string, sentAt *time.Time) error
}

type NotificationProvider interface {
    SendEmail(to, subject, content string) error
    SendSMS(to, content string) error
}

type NotificationService interface {
    SendNotification(userID uuid.UUID, notificationType, channel, subject, content string, metadata map[string]interface{}) error
    GetNotifications(userID uuid.UUID, limit, offset int) ([]*Notification, error)
    MarkAsRead(id string) error
}
