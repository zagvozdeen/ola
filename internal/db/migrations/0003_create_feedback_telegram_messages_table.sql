-- +goose up
CREATE TABLE IF NOT EXISTS feedback_telegram_messages
(
    feedback_id INTEGER REFERENCES feedback (id) ON DELETE CASCADE NOT NULL,
    chat_id     BIGINT                                             NOT NULL,
    message_id  BIGINT                                             NOT NULL,
    PRIMARY KEY (feedback_id, chat_id, message_id)
);

-- +goose down
DROP TABLE IF EXISTS feedback_telegram_messages;
