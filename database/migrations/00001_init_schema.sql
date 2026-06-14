-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `users` (
    `id`                   VARCHAR(36)  NOT NULL,
    `name`                 VARCHAR(100) NOT NULL,
    `email`                VARCHAR(100) NOT NULL,
    `password_hash`        VARCHAR(255) NOT NULL,
    `role`                 VARCHAR(20)  NOT NULL,
    `department_id`        VARCHAR(36)  NULL,
    `creator_id`           VARCHAR(36)  NULL,
    `fixed_ip`             VARCHAR(45)  NULL,
    `subnet`               VARCHAR(45)  NULL,
    `created_at`           DATETIME(3)  NULL,
    `updated_at`           DATETIME(3)  NULL,
    `is_online`            TINYINT(1)   NOT NULL DEFAULT 0,
    `last_connection_time` DATETIME(3)  NULL,
    `real_address`         VARCHAR(45)  NULL,
    `virtual_address`      VARCHAR(45)  NULL,
    `bytes_received`       BIGINT       NOT NULL DEFAULT 0,
    `bytes_sent`           BIGINT       NOT NULL DEFAULT 0,
    `connected_since`      DATETIME(3)  NULL,
    `last_ref`             DATETIME(3)  NULL,
    `online_duration`      BIGINT       NOT NULL DEFAULT 0,
    `is_paused`            TINYINT(1)   NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_users_name` (`name`),
    UNIQUE KEY `uni_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `departments` (
    `id`         VARCHAR(36)  NOT NULL,
    `name`       VARCHAR(100) NOT NULL,
    `head_id`    VARCHAR(36)  NULL,
    `parent_id`  VARCHAR(36)  NULL,
    `created_at` DATETIME(3)  NULL,
    `updated_at` DATETIME(3)  NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_departments_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `client_logs` (
    `id`                   VARCHAR(36)  NOT NULL,
    `user_id`              VARCHAR(36)  NULL,
    `is_online`            TINYINT(1)   NOT NULL DEFAULT 0,
    `real_address`         VARCHAR(255) NULL,
    `online_duration`      BIGINT       NOT NULL DEFAULT 0,
    `traffic_usage`        BIGINT       NOT NULL DEFAULT 0,
    `last_connection_time` DATETIME(3)  NULL,
    `created_at`           DATETIME(3)  NULL,
    PRIMARY KEY (`id`),
    KEY `idx_client_logs_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `client_logs`;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS `departments`;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS `users`;
-- +goose StatementEnd
