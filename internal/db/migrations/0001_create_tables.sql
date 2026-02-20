-- +goose up
CREATE TABLE IF NOT EXISTS stock
(
    sku       TEXT PRIMARY KEY,
    available INT NOT NULL
);

CREATE TABLE IF NOT EXISTS reservations
(
    reservation_id UUID PRIMARY KEY,
    saga_id        UUID        NOT NULL,
    sku            TEXT        NOT NULL,
    qty            INTEGER     NOT NULL,
    status         TEXT        NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS idempotency
(
    idempotency_key UUID PRIMARY KEY,
    status          TEXT        NOT NULL,
    response        JSONB       NULL,
    created_at      TIMESTAMPTZ NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS stock;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS idempotency;
