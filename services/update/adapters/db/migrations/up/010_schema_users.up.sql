CREATE TABLE users
(
    user_id       BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname      VARCHAR(64),
    avatar_url    VARCHAR(255),
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    is_admin      BOOLEAN   NOT NULL DEFAULT false
);