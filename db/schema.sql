create table users
(
    id            serial primary key,
    username      text unique not null,
    password_hash text        not null
);

CREATE TABLE chats
(
    id        SERIAL PRIMARY KEY,
    name      TEXT    NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    owner_id  INT REFERENCES users (id)
);

CREATE TABLE chat_user
(
    chat_id INT REFERENCES chats (id),
    user_id INT REFERENCES users (id),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages
(
    id         SERIAL PRIMARY KEY,
    chat_id    INT REFERENCES chats (id) NOT NULL,
    user_id    INT REFERENCES users (id) NOT NULL ,
    text       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Outbox для сообщений
CREATE TABLE IF NOT EXISTS message_outbox (
    id              BIGSERIAL PRIMARY KEY,
    message_id      INT UNIQUE REFERENCES messages (id),
    payload         JSONB NOT NULL,
    published_at    TIMESTAMP NULL,
    attempts        INT NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION trg_messages_after_insert()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO message_outbox (message_id, payload)
    VALUES (NEW.id, jsonb_build_object(
        'msg_type', 'evt.message.posted',
        'message_id', NEW.id,
        'chat_id', NEW.chat_id,
        'sender_id', NEW.user_id,
        'text', NEW.text,
        'created_unix', EXTRACT(EPOCH FROM NEW.created_at)
    ));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS messages_after_insert ON messages;
CREATE TRIGGER messages_after_insert
AFTER INSERT ON messages
FOR EACH ROW EXECUTE PROCEDURE trg_messages_after_insert();

-- Saga orchestrator storage
CREATE TABLE IF NOT EXISTS saga_instances (
    id           BIGSERIAL PRIMARY KEY,
    saga_id      TEXT UNIQUE NOT NULL,
    type         TEXT NOT NULL,
    status       TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS saga_step_log (
    id           BIGSERIAL PRIMARY KEY,
    saga_id      TEXT NOT NULL,
    step         INT  NOT NULL,
    action       TEXT NOT NULL,
    status       TEXT NOT NULL,
    payload      JSONB,
    error        TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Simple DLQ table (optional backup for failed msgs)
CREATE TABLE IF NOT EXISTS dlq_messages (
    id           BIGSERIAL PRIMARY KEY,
    topic        TEXT NOT NULL,
    key          BYTEA,
    value        BYTEA,
    headers      JSONB,
    error        TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
