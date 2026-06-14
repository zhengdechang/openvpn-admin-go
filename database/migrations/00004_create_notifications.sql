-- +goose Up
CREATE TABLE IF NOT EXISTS notifications (
    id         VARCHAR(36)  PRIMARY KEY,
    type       VARCHAR(50)  NOT NULL,
    user_name  VARCHAR(100) NOT NULL,
    real_ip    VARCHAR(45)  DEFAULT '',
    virtual_ip VARCHAR(45)  DEFAULT '',
    is_read    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read    ON notifications (is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_user_name  ON notifications (user_name);

-- +goose Down
DROP TABLE IF EXISTS notifications;
