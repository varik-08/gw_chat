package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    kgo "github.com/segmentio/kafka-go"
    "github.com/prometheus/client_golang/prometheus/promhttp"

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
    db := conf.InitDB(cfg.DB)

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
                // log receive
                _, _ = db.Exec(ctx, `INSERT INTO saga_instances (saga_id, type, status) 
                    VALUES ($1,$2,$3) ON CONFLICT (saga_id) DO UPDATE SET updated_at=NOW()`,
                    cmd.SagaID, "generic", "in_progress")
                _, _ = db.Exec(ctx, `INSERT INTO saga_step_log (saga_id, step, action, status, payload) VALUES ($1,$2,$3,$4,$5)`,
                    cmd.SagaID, cmd.Step, m.Topic, "received", m.Value)
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
                if err := prod.Publish(ctx, evtTopic, key, m.Value, headers); err != nil {
                    _, _ = db.Exec(ctx, `INSERT INTO saga_step_log (saga_id, step, action, status, payload, error) VALUES ($1,$2,$3,$4,$5,$6)`,
                        cmd.SagaID, cmd.Step, evtTopic, "publish_error", m.Value, err.Error())
                    return err
                }
                _, _ = db.Exec(ctx, `INSERT INTO saga_step_log (saga_id, step, action, status, payload) VALUES ($1,$2,$3,$4,$5)`,
                    cmd.SagaID, cmd.Step, evtTopic, "published", m.Value)
                _, _ = db.Exec(ctx, `UPDATE saga_instances SET status=$2, updated_at=NOW() WHERE saga_id=$1`, cmd.SagaID, "completed")
                return nil
            })
        }(cons)
    }

    // health/metrics
    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
        mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
        mux.Handle("/metrics", promhttp.Handler())
        _ = http.ListenAndServe(":9103-health", mux)
    }()

    // Блокируемся
    for { time.Sleep(10 * time.Second) }
}

