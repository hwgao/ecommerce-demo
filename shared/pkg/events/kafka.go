package events

import (
    "encoding/json"
    "time"
    "github.com/IBM/sarama"
    "github.com/google/uuid"
)

type EventBus interface {
    Publish(topic string, event interface{}) error
    Subscribe(topic string, handler func([]byte) error) error
}

type KafkaEventBus struct {
    producer sarama.SyncProducer
    consumer sarama.Consumer
}

func NewKafkaEventBus(brokers string) *KafkaEventBus {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    
    producer, _ := sarama.NewSyncProducer([]string{brokers}, config)
    consumer, _ := sarama.NewConsumer([]string{brokers}, config)
    
    return &KafkaEventBus{
        producer: producer,
        consumer: consumer,
    }
}

func (k *KafkaEventBus) Publish(topic string, event interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(data),
    }
    
    _, _, err = k.producer.SendMessage(msg)
    return err
}

func (k *KafkaEventBus) Subscribe(topic string, handler func([]byte) error) error {
    partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        return err
    }
    
    go func() {
        for message := range partitionConsumer.Messages() {
            handler(message.Value)
        }
    }()
    
    return nil
}

// Event types
type UserRegisteredEvent struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Timestamp time.Time `json:"timestamp"`
}

type OrderCreatedEvent struct {
    OrderID     uuid.UUID       `json:"order_id"`
    UserID      uuid.UUID       `json:"user_id"`
    TotalAmount float64         `json:"total_amount"`
    Currency    string          `json:"currency"`
    Items       []interface{}   `json:"items"`
    Timestamp   time.Time       `json:"timestamp"`
}

type OrderStatusUpdatedEvent struct {
    OrderID   uuid.UUID `json:"order_id"`
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
}
