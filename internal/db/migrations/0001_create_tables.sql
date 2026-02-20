-- +goose up
CREATE TYPE user_role AS ENUM ('user', 'moderator', 'admin');

CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    tid        BIGINT       NULL UNIQUE,
    uuid       UUID         NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NULL,
    username   VARCHAR(255) NULL,
    email      VARCHAR(256) NULL,
    password   VARCHAR(256) NULL,
    role       user_role    NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL
);

CREATE TABLE IF NOT EXISTS files
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID                          NOT NULL UNIQUE,
    content     VARCHAR(255)                  NOT NULL,
    size        INTEGER                       NOT NULL,
    mime_type   VARCHAR(255)                  NOT NULL,
    origin_name VARCHAR(255)                  NOT NULL,
    user_id     INTEGER REFERENCES users (id) NOT NULL,
    created_at  TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS products
(
    id         UUID PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) NOT NULL,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS services
(
    id         UUID PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) NOT NULL,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS orders
(
    id         UUID PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) NOT NULL,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS feedback
(
    id         UUID PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) NOT NULL,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL
);

-- +goose down
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
