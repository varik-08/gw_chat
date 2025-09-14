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
    reader  *kgo.Reader
    brokers []string
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
    return &Consumer{reader: r, brokers: brokers}
}

// Start запускает цикл чтения. Идемпотентность обязан обеспечить handler.
func (c *Consumer) Start(ctx context.Context, handler Handler) error {
    // Локальный продюсер для DLQ. В реальном коде лучше инжектить.
    dlqWriter := &kgo.Writer{Addr: kgo.TCP(c.brokers...), AllowAutoTopicCreation: true}
    defer func() { _ = dlqWriter.Close() }()

    for {
        m, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if errors.Is(err, context.Canceled) {
                return nil
            }
            return err
        }

        var ok bool
        var lastErr error
        backoff := 100 * time.Millisecond
        for attempt := 0; attempt < 5; attempt++ {
            if err := handler(ctx, m); err != nil {
                lastErr = err
                time.Sleep(backoff)
                backoff *= 2
                continue
            }
            ok = true
            break
        }

        if ok {
            if err := c.reader.CommitMessages(ctx, m); err != nil {
                return err
            }
            continue
        }

        // Отправляем в DLQ
        _ = dlqWriter.WriteMessages(ctx, kgo.Message{
            Topic:   m.Topic + ".dlq",
            Key:     m.Key,
            Value:   m.Value,
            Headers: m.Headers,
            Time:    time.Now(),
        })
        // Коммитим, чтобы не застревать на ядовитом сообщении
        if err := c.reader.CommitMessages(ctx, m); err != nil {
            return err
        }
        _ = lastErr
    }
}

// Close закрывает reader.
func (c *Consumer) Close() error {
    if c == nil || c.reader == nil {
        return nil
    }
    return c.reader.Close()
}

