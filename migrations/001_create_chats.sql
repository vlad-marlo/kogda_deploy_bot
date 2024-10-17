CREATE TYPE chat_type AS ENUM (
    'private',
    'group',
    'supergroup',
    'channel',
    'privatechannel'
    );
CREATE TABLE chats
(
    id         BIGINT PRIMARY KEY NOT NULL,
    chat_type  chat_type          NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
---- create above / drop below ----

DROP TABLE chats;
DROP TYPE chat_type;