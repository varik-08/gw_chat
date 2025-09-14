package main

import (
    "context"
    "encoding/json"
    "log"
    "os"

    kgo "github.com/segmentio/kafka-go"

    conf "github.com/varik-08/gw_chat/config"
    mqk "github.com/varik-08/gw_chat/internal/mq/kafka"
    otelinit "github.com/varik-08/gw_chat/internal/otel"
)

type sagaCmd struct {
    Type   string          `json:"type"`
    SagaID string          `json:"saga_id"`
    Step   int             `json:"step"`
    Data   json.RawMessage `json:"data"`
}

func main() {
    ctx := context.Background()
    shutdown, err := otelinit.Init(ctx, "saga-orchestrator")
    if err != nil { log.Fatalf("otel: %v", err) }
    defer func() { _ = shutdown(context.Background()) }()

    cfg, err := conf.GetConfig()
    if err != nil { log.Fatalf("config: %v", err) }
    _ = conf.InitDB(cfg.DB)

    brokers := os.Getenv("KAFKA_BROKERS")
    if brokers == "" { brokers = "kafka:9092" }
    prod := mqk.NewProducer([]string{brokers}, "orchestrator")

    // Подписываемся на команды
    topics := []string{"cmd.chat.create", "cmd.chat.add_members", "cmd.message.post"}
    for _, t := range topics {
        cons := mqk.NewConsumer([]string{brokers}, "orchestrator", t)
        go func(c *mqk.Consumer) {
            _ = c.Start(ctx, func(ctx context.Context, m kgo.Message) error {
                var cmd sagaCmd
                _ = json.Unmarshal(m.Value, &cmd)
                // На минималках: echo -> evt.*
                evtTopic := map[string]string{
                    "cmd.chat.create":     "evt.chat.created",
                    "cmd.chat.add_members": "evt.chat.members_added",
                    "cmd.message.post":     "evt.message.posted",
                }[m.Topic]
                if evtTopic == "" { return nil }
                headers := map[string]string{"msg_type": evtTopic}
                key := m.Key
                if len(key) == 0 && cmd.SagaID != "" { key = []byte(cmd.SagaID) }
                return prod.Publish(ctx, evtTopic, key, m.Value, headers)
            })
        }(cons)
    }

    // Блокируемся
    select {}
}

