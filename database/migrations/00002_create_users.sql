-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id                   VARCHAR(36)  PRIMARY KEY,
    name                 VARCHAR(100) NOT NULL,
    email                VARCHAR(100) NOT NULL,
    password_hash        VARCHAR(255) NOT NULL,
    role                 VARCHAR(20)  NOT NULL DEFAULT 'user',
    department_id        VARCHAR(36)  NOT NULL DEFAULT '',
    creator_id           VARCHAR(36)  NOT NULL DEFAULT '',
    fixed_ip             VARCHAR(45)  NOT NULL DEFAULT '',
    subnet               VARCHAR(45)  NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    is_online            BOOLEAN      NOT NULL DEFAULT FALSE,
    last_connection_time TIMESTAMPTZ,
    real_address         VARCHAR(45)  NOT NULL DEFAULT '',
    virtual_address      VARCHAR(45)  NOT NULL DEFAULT '',
    bytes_received       BIGINT       NOT NULL DEFAULT 0,
    bytes_sent           BIGINT       NOT NULL DEFAULT 0,
    connected_since      TIMESTAMPTZ,
    last_ref             TIMESTAMPTZ,
    online_duration      BIGINT       NOT NULL DEFAULT 0,
    is_paused            BOOLEAN      NOT NULL DEFAULT FALSE,
    CONSTRAINT users_name_key  UNIQUE (name),
    CONSTRAINT users_email_key UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
