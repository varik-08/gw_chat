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
)
