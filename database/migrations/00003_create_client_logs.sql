-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS client_logs (
    id                   VARCHAR(36)  PRIMARY KEY,
    user_id              VARCHAR(36)  NOT NULL DEFAULT '',
    is_online            BOOLEAN      NOT NULL DEFAULT FALSE,
    real_address         VARCHAR(255) NOT NULL DEFAULT '',
    online_duration      BIGINT       NOT NULL DEFAULT 0,
    traffic_usage        BIGINT       NOT NULL DEFAULT 0,
    last_connection_time TIMESTAMPTZ,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_client_logs_user_id ON client_logs(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_client_logs_user_id;
DROP TABLE IF EXISTS client_logs;
-- +goose StatementEnd
