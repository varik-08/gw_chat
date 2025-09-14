package kafka

import (
    "context"
    "time"

    kgo "github.com/segmentio/kafka-go"
)

// Producer инкапсулирует синхронную публикацию сообщений в Kafka.
type Producer struct {
    writer *kgo.Writer
}

// NewProducer создаёт продюсер с явной передачей брокеров.
func NewProducer(brokers []string, clientID string) *Producer {
    w := &kgo.Writer{
        Addr:                   kgo.TCP(brokers...),
        Balancer:               &kgo.Hash{},
        RequiredAcks:           kgo.RequireAll,
        Async:                  false,
        AllowAutoTopicCreation: true,
        BatchTimeout:           10 * time.Millisecond,
    }
    _ = clientID // зарезервировано для будущего трейсинга/метрик
    return &Producer{writer: w}
}

// Publish отправляет одно сообщение в указанный топик.
func (p *Producer) Publish(ctx context.Context, topic string, key []byte, value []byte, headers map[string]string) error {
    if p == nil || p.writer == nil {
        return nil
    }

    var kh []kgo.Header
    for k, v := range headers {
        kh = append(kh, kgo.Header{Key: k, Value: []byte(v)})
    }

    msg := kgo.Message{
        Topic:   topic,
        Key:     key,
        Value:   value,
        Headers: kh,
        Time:    time.Now(),
    }

    return p.writer.WriteMessages(ctx, msg)
}

// Close закрывает внутренний writer.
func (p *Producer) Close() error {
    if p == nil || p.writer == nil {
        return nil
    }
    return p.writer.Close()
}

