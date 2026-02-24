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

CREATE TABLE IF NOT EXISTS products
(
    id          SERIAL PRIMARY KEY,
    uuid        UUID                          NOT NULL UNIQUE,
    name        VARCHAR(255)                  NOT NULL,
    description TEXT                          NOT NULL,
    price_from  INTEGER                       NOT NULL,
    price_to    INTEGER                       NULL,
    type        product_type                  NOT NULL,
    file_id     INTEGER REFERENCES files (id) NOT NULL,
    user_id     INTEGER REFERENCES users (id) NOT NULL,
    created_at  TIMESTAMPTZ                   NOT NULL,
    updated_at  TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS category_product
(
    category_id INTEGER REFERENCES categories (id) NOT NULL,
    product_id  INTEGER REFERENCES products (id)   NOT NULL,
    PRIMARY KEY (category_id, product_id)
);

CREATE TABLE IF NOT EXISTS reviews
(
    id           SERIAL PRIMARY KEY,
    uuid         UUID                          NOT NULL UNIQUE,
    name         VARCHAR(255)                  NOT NULL,
    content      TEXT                          NOT NULL,
    file_id      INTEGER REFERENCES files (id) NOT NULL,
    user_id      INTEGER REFERENCES users (id) NOT NULL,
    published_at TIMESTAMPTZ                   NOT NULL,
    created_at   TIMESTAMPTZ                   NOT NULL,
    updated_at   TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS carts
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                          NOT NULL UNIQUE,
    user_id    INTEGER REFERENCES users (id) NULL UNIQUE,
    session_id UUID                          NULL UNIQUE,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL,
    CHECK (user_id IS NOT NULL OR session_id IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS cart_items
(
    cart_id    INTEGER REFERENCES carts (id) ON DELETE CASCADE NOT NULL,
    product_id INTEGER REFERENCES products (id)                 NOT NULL,
    qty        INTEGER                                          NOT NULL CHECK (qty > 0),
    PRIMARY KEY (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                          NOT NULL UNIQUE,
    source     order_source                  NOT NULL,
    name       VARCHAR(255)                  NOT NULL,
    phone      VARCHAR(255)                  NOT NULL,
    content    TEXT                          NOT NULL,
    user_id    INTEGER REFERENCES users (id) NULL,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items
(
    order_id      INTEGER REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    product_id    INTEGER REFERENCES products (id)                 NOT NULL,
    product_name  VARCHAR(255)                                     NOT NULL,
    price_from    INTEGER                                          NOT NULL,
    price_to      INTEGER                                          NULL,
    qty           INTEGER                                          NOT NULL CHECK (qty > 0),
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS feedback
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                          NOT NULL UNIQUE,
    name       VARCHAR(255)                  NOT NULL,
    phone      VARCHAR(255)                  NOT NULL,
    content    TEXT                          NOT NULL,
    user_id    INTEGER REFERENCES users (id) NULL,
    created_at TIMESTAMPTZ                   NOT NULL,
    updated_at TIMESTAMPTZ                   NOT NULL
);

-- DEV ONLY (DELETE IN PRODUCTION)
INSERT INTO users (tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at)
VALUES (NULL, gen_random_uuid(), 'Ivan', 'Ivanov', 'ivan', 'ivan@mail.ru', '$2a$10$5IGmMleZEgX8azGohTVBaeAlSnONTuRLVT0YxtP6HXzsQOrabJvbu', 'admin', NOW(), NOW());

INSERT INTO files (uuid, content, size, mime_type, origin_name, user_id, created_at)
VALUES (gen_random_uuid(), '/files/1.jpg', 0, 'image/jpeg', '1.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/2.jpg', 0, 'image/jpeg', '2.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/3.jpg', 0, 'image/jpeg', '3.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/4.jpg', 0, 'image/jpeg', '4.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/5.jpg', 0, 'image/jpeg', '5.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/6.jpg', 0, 'image/jpeg', '6.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/7.jpg', 0, 'image/jpeg', '7.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/8.jpg', 0, 'image/jpeg', '8.jpg', 1, NOW()),
       (gen_random_uuid(), '/files/9.jpg', 0, 'image/jpeg', '9.jpg', 1, NOW());

INSERT INTO products (uuid, name, description, price_from, price_to, type, file_id, user_id, created_at, updated_at)
VALUES (gen_random_uuid(), 'Фонтан из воздушных шаров', 'Композиция по индивидуальному дизайну для любого события',
        3500,
        NULL, 'product', 1, 1,
        NOW(), NOW()),
       (gen_random_uuid(), 'Оформление помещения / фотозона',
        'Декорирование любого помещения по индивидуальному дизайну',
        7000, NULL, 'product',
        2, 1, NOW(), NOW()),
       (gen_random_uuid(), 'Коробка - сюрприз', 'Подарочный бокс с композицией из шаров для любого события', 5000,
        NULL, 'product', 3, 1,
        NOW(),
        NOW()),
       (gen_random_uuid(), 'Бабл бокс', 'Креативная упаковка для небольшого подарка с шаром баблс', 3000, NULL, 'product', 4, 1,
        NOW(),
        NOW()),
       (gen_random_uuid(), 'Фонтан из воздушных шаров', 'Композиция по индивидуальному дизайну для любого события',
        3500,
        NULL, 'service', 1, 1,
        NOW(), NOW()),
       (gen_random_uuid(), 'Оформление помещения / фотозона',
        'Декорирование любого помещения по индивидуальному дизайну',
        7000, NULL, 'service',
        2, 1, NOW(), NOW()),
       (gen_random_uuid(), 'Коробка - сюрприз', 'Подарочный бокс с композицией из шаров для любого события', 5000,
        NULL, 'service', 3, 1,
        NOW(),
        NOW()),
       (gen_random_uuid(), 'Бабл бокс', 'Креативная упаковка для небольшого подарка с шаром баблс', 3000, NULL, 'service', 4, 1,
        NOW(),
        NOW());

INSERT INTO categories (uuid, name, created_at, updated_at)
VALUES (gen_random_uuid(), 'Детские праздники', NOW(), NOW()),
       (gen_random_uuid(), 'Корпоратив', NOW(), NOW()),
       (gen_random_uuid(), 'День рождение', NOW(), NOW()),
       (gen_random_uuid(), 'Свадьба', NOW(), NOW()),
       (gen_random_uuid(), 'Выписка', NOW(), NOW()),
       (gen_random_uuid(), 'Гендер пати', NOW(), NOW()),
       (gen_random_uuid(), '8 марта', NOW(), NOW()),
       (gen_random_uuid(), '23 февраля', NOW(), NOW()),
       (gen_random_uuid(), 'Без повода', NOW(), NOW());

INSERT INTO reviews (uuid, name, content, file_id, user_id, published_at, created_at, updated_at)
VALUES (gen_random_uuid(), 'Елена',
        'Огромное спасибо за шарики! Именинник был в восторге и доставка порадовала, все вовремя) Не пожалела, что обратилась именно к вам!',
        5, 1, NOW(), NOW(), NOW()),
       (gen_random_uuid(), 'Евгений',
        'Обратился за бабл боксом с украшением внутри, хотелось поздравить девушку с годовщиной. Она очень удивилась такой креативной идее. Я тоже не видел ничего подобного в нашем городе до этого)) Желаю вам дальнейшего развития, вы классные!',
        6, 1, NOW(), NOW(), NOW()),
       (gen_random_uuid(), 'Александр',
        'Выражаю благодарность студии за профессиональное оформление нашего корпоратива. Требовался лаконичный декор в цветах компании. Результат превзошел ожидания: композиции у входа и фотозона были выполнены безупречно, с вниманием к деталям. Отдельно отмечу пунктуальность, четкое соблюдение сроков и договоренностей.',
        7, 1, NOW(), NOW(), NOW()),
       (gen_random_uuid(), 'Любовь',
        'Благодарю персонал за то, что взялись за очень срочный заказ, выполнили и доставили максимально быстро. Посоветую вас друзьям и сама обращусь еще не раз.',
        8, 1, NOW(), NOW(), NOW()),
       (gen_random_uuid(), 'Анна',
        'Как человек, который ценит визуал, долго искала в Екатеринбурге студию, которая умеет в тренды. Команда Ola предложила крутые цветовые сочетания для моей вечеринки. Шары держались несколько дней, не теряя вид! Ни одного лопнувшего. Это показатель. Если вам важен дизайн, атмосфера и стойкость - выбор очевиден.',
        9, 1, NOW(), NOW(), NOW());

-- +goose down
DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS category_product;
DROP TABLE IF EXISTS products;
DROP TYPE IF EXISTS product_type;
DROP TYPE IF EXISTS order_source;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
