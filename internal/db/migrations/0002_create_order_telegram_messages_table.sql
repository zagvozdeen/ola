-- +goose up
CREATE TABLE IF NOT EXISTS order_telegram_messages
(
    order_id   INTEGER REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    chat_id    BIGINT                                           NOT NULL,
    message_id BIGINT                                           NOT NULL,
    PRIMARY KEY (order_id, chat_id, message_id)
);

-- +goose down
DROP TABLE IF EXISTS order_telegram_messages;
