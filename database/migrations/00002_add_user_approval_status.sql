-- +goose Up
-- +goose StatementBegin
ALTER TABLE `users` ADD COLUMN `approval_status` VARCHAR(20) NOT NULL DEFAULT 'approved';
-- +goose StatementEnd
-- +goose StatementBegin
UPDATE `users` SET `approval_status` = 'approved';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `users` DROP COLUMN `approval_status`;
-- +goose StatementEnd
