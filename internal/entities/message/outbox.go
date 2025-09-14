package message

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    mqk "github.com/varik-08/gw_chat/internal/mq/kafka"
)

type OutboxPublisher struct {
    db       *pgxpool.Pool
    producer *mqk.Producer
}

func NewOutboxPublisher(db *pgxpool.Pool, producer *mqk.Producer) *OutboxPublisher {
    return &OutboxPublisher{db: db, producer: producer}
}

func (p *OutboxPublisher) Run(ctx context.Context) error {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            if err := p.tick(ctx); err != nil {
                // лог-с IL
            }
        }
    }
}

func (p *OutboxPublisher) tick(ctx context.Context) error {
    rows, err := p.db.Query(ctx, `
        SELECT id, payload FROM message_outbox
        WHERE published_at IS NULL AND next_attempt_at <= NOW()
        ORDER BY id ASC
        LIMIT 50
    `)
    if err != nil { return err }
    defer rows.Close()

    type rec struct { id int64; payload []byte }
    var batch []rec
    for rows.Next() {
        var id int64
        var payload []byte
        if err := rows.Scan(&id, &payload); err != nil { return err }
        batch = append(batch, rec{id: id, payload: payload})
    }
    if len(batch) == 0 { return nil }

    for _, r := range batch {
        var env map[string]any
        _ = json.Unmarshal(r.payload, &env)
        msgType, _ := env["msg_type"].(string)
        topic := msgType
        if topic == "" { topic = "evt.message.posted" }
        key := []byte{}
        if v, ok := env["chat_id"].(float64); ok { key = []byte(fmt.Sprintf("%d", int64(v))) }
        if err := p.producer.Publish(ctx, topic, key, r.payload, map[string]string{"msg_type": topic}); err != nil {
            // backoff
            _, _ = p.db.Exec(ctx, `UPDATE message_outbox SET attempts = attempts + 1, next_attempt_at = NOW() + INTERVAL '1 second' * (2 ^ LEAST(attempts,10)) WHERE id=$1`, r.id)
            continue
        }
        _, _ = p.db.Exec(ctx, `UPDATE message_outbox SET published_at = NOW() WHERE id=$1`, r.id)
    }

    return nil
}