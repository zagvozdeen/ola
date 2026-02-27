-- +goose up
CREATE TYPE user_role AS ENUM ('user', 'manager', 'moderator', 'admin');

CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    tid        BIGINT       NULL UNIQUE,
    uuid       UUID         NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NULL,
    username   VARCHAR(255) NULL,
    email      VARCHAR(256) NULL,
    phone      VARCHAR(255) NULL,
    password   VARCHAR(256) NULL,
    role       user_role    NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL
);

CREATE TABLE IF NOT EXISTS files
(
--     id          SERIAL PRIMARY KEY,
    uuid        UUID                          NOT NULL UNIQUE,
    content     VARCHAR(255)                  NOT NULL UNIQUE PRIMARY KEY,
    size        BIGINT                        NOT NULL,
    mime_type   VARCHAR(255)                  NOT NULL,
    origin_name VARCHAR(255)                  NOT NULL,
    user_id     INTEGER REFERENCES users (id) NOT NULL,
    created_at  TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS categories
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID         NOT NULL UNIQUE,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL
);

CREATE TYPE product_type AS ENUM ('product', 'service');
CREATE TYPE order_source AS ENUM ('landing', 'spa', 'tma');
CREATE TYPE request_status AS ENUM ('created', 'in_progress', 'reviewed');

CREATE TABLE IF NOT EXISTS products
(
    id           SERIAL PRIMARY KEY,
    uuid         UUID                                                  NOT NULL UNIQUE,
    name         VARCHAR(255)                                          NOT NULL,
    description  TEXT                                                  NOT NULL,
    price_from   INTEGER                                               NOT NULL,
    price_to     INTEGER                                               NULL,
    type         product_type                                          NOT NULL,
    is_main      BOOLEAN                                               NOT NULL,
    file_content INTEGER REFERENCES files (content) ON DELETE RESTRICT NOT NULL,
    user_id      INTEGER REFERENCES users (id) ON DELETE RESTRICT      NOT NULL,
    created_at   TIMESTAMPTZ                                           NOT NULL,
    updated_at   TIMESTAMPTZ                                           NOT NULL
);

CREATE TABLE IF NOT EXISTS category_product
(
    category_id INTEGER REFERENCES categories (id) NOT NULL,
    product_id  INTEGER REFERENCES products (id)   NOT NULL,
    PRIMARY KEY (category_id, product_id)
);

CREATE TABLE IF NOT EXISTS carts
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                            NOT NULL UNIQUE,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE NULL UNIQUE,
    session_id UUID                                            NULL UNIQUE,
    created_at TIMESTAMPTZ                                     NOT NULL,
    updated_at TIMESTAMPTZ                                     NOT NULL,
    CHECK (user_id IS NOT NULL OR session_id IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS cart_items
(
    cart_id    INTEGER REFERENCES carts (id) ON DELETE CASCADE    NOT NULL,
    product_id INTEGER REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    qty        INTEGER                                            NOT NULL CHECK (qty > 0),
    PRIMARY KEY (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                             NOT NULL UNIQUE,
    status     request_status                                   NOT NULL,
    source     order_source                                     NOT NULL,
    name       VARCHAR(255)                                     NOT NULL,
    phone      VARCHAR(255)                                     NOT NULL,
    content    TEXT                                             NOT NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE RESTRICT NULL,
    created_at TIMESTAMPTZ                                      NOT NULL,
    updated_at TIMESTAMPTZ                                      NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items
(
    order_id     INTEGER REFERENCES orders (id) ON DELETE RESTRICT   NOT NULL,
    product_id   INTEGER REFERENCES products (id) ON DELETE RESTRICT NOT NULL,
    product_name VARCHAR(255)                                        NOT NULL,
    price_from   INTEGER                                             NOT NULL,
    price_to     INTEGER                                             NULL,
    qty          INTEGER                                             NOT NULL CHECK (qty > 0),
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS order_comments
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                              NOT NULL UNIQUE,
    content    TEXT                                              NOT NULL,
    order_id   INTEGER REFERENCES orders (id) ON DELETE RESTRICT NOT NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE RESTRICT  NOT NULL,
    created_at TIMESTAMPTZ                                       NOT NULL,
    updated_at TIMESTAMPTZ                                       NOT NULL
);

CREATE TABLE IF NOT EXISTS feedback
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                             NOT NULL UNIQUE,
    status     request_status                                   NOT NULL,
    source     order_source                                     NOT NULL,
    type       VARCHAR(64)                                      NOT NULL,
    name       VARCHAR(255)                                     NOT NULL,
    phone      VARCHAR(255)                                     NOT NULL,
    content    TEXT                                             NOT NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE RESTRICT NOT NULL,
    created_at TIMESTAMPTZ                                      NOT NULL,
    updated_at TIMESTAMPTZ                                      NOT NULL
);

CREATE TABLE IF NOT EXISTS order_telegram_messages
(
    order_id   INTEGER REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    chat_id    BIGINT                                           NOT NULL,
    message_id BIGINT                                           NOT NULL,
    PRIMARY KEY (order_id, chat_id, message_id)
);

CREATE TABLE IF NOT EXISTS feedback_telegram_messages
(
    feedback_id INTEGER REFERENCES feedback (id) ON DELETE CASCADE NOT NULL,
    chat_id     BIGINT                                             NOT NULL,
    message_id  BIGINT                                             NOT NULL,
    PRIMARY KEY (feedback_id, chat_id, message_id)
);

CREATE TABLE IF NOT EXISTS actions
(
    id         SERIAL PRIMARY KEY,
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS feedback_telegram_messages;
DROP TABLE IF EXISTS order_telegram_messages;
DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS order_comments;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS category_product;
DROP TABLE IF EXISTS products;
DROP TYPE IF EXISTS product_type;
DROP TYPE IF EXISTS order_source;
DROP TYPE IF EXISTS request_status;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
