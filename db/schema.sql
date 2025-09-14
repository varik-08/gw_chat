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
