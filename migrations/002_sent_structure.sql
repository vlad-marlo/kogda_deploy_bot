create table chat_day
(
    chat_id BIGINT REFERENCES chats (id),
    day     DATE,
    message TEXT,
    CONSTRAINT chat_day_pk PRIMARY KEY (chat_id, day, message)
);
---- create above / drop below ----
drop table chat_day;