package kafka

import (
    "context"
    "errors"
    "time"

    kgo "github.com/segmentio/kafka-go"
)

// Handler — пользовательская функция обработки сообщений.
type Handler func(ctx context.Context, msg kgo.Message) error

// Consumer инкапсулирует чтение из Kafka c коммитом оффсетов при успехе.
type Consumer struct {
    reader *kgo.Reader
}

// NewConsumer создаёт консьюмера группы.
func NewConsumer(brokers []string, groupID, topic string) *Consumer {
    r := kgo.NewReader(kgo.ReaderConfig{
        Brokers:        brokers,
        GroupID:        groupID,
        Topic:          topic,
        MinBytes:       1,
        MaxBytes:       10 * 1024 * 1024,
        CommitInterval: 0, // управляем коммитом вручную
    })
    return &Consumer{reader: r}
}

// Start запускает цикл чтения. Идемпотентность обязан обеспечить handler.
func (c *Consumer) Start(ctx context.Context, handler Handler) error {
    for {
        m, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if errors.Is(err, context.Canceled) {
                return nil
            }
            return err
        }

        if err := handler(ctx, m); err != nil {
            // Не коммитим — сообщение будет переобработано (at-least-once)
            // Небольшая пауза, чтобы избежать tight loop
            time.Sleep(200 * time.Millisecond)
            continue
        }
        if err := c.reader.CommitMessages(ctx, m); err != nil {
            return err
        }
    }
}

// Close закрывает reader.
func (c *Consumer) Close() error {
    if c == nil || c.reader == nil {
        return nil
    }
    return c.reader.Close()
}

