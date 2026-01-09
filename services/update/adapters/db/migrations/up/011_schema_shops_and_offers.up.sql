CREATE TABLE shops
(
    shop_id BIGSERIAL PRIMARY KEY,
    name    VARCHAR(128) NOT NULL,
    url     VARCHAR(255)
);

CREATE TABLE product_offers
(
    offer_id       BIGSERIAL PRIMARY KEY,
    shop_id        BIGINT         NOT NULL REFERENCES shops (shop_id) ON DELETE CASCADE,
    component_type VARCHAR(32)    NOT NULL, -- 'CPU','GPU','MOTHERBOARD',...
    component_id   BIGINT         NOT NULL,
    price          NUMERIC(12, 2) NOT NULL,
    available      BOOLEAN        NOT NULL DEFAULT TRUE
    -- Полиморфный component_id.
    -- Ссылки на конкретные таблицы будем контролировать на уровне бизнес-логики/триггеров.
);