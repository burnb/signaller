-- +goose Up
-- +goose StatementBegin
ALTER TABLE positions
    MODIFY COLUMN created_at BIGINT UNSIGNED NOT NULL,
    MODIFY COLUMN updated_at BIGINT UNSIGNED NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE positions
    MODIFY COLUMN updated_at TIMESTAMP NOT NULL,
    MODIFY COLUMN created_at TIMESTAMP NOT NULL;
-- +goose StatementEnd
