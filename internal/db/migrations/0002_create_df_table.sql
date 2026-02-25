-- +goose up
CREATE TABLE IF NOT EXISTS order_telegram_messages
(
    order_id   INTEGER REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    chat_id    BIGINT                                           NOT NULL,
    message_id BIGINT                                           NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS order_telegram_messages;
