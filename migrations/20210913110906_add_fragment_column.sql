-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls ADD COLUMN fragment VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls DROP COLUMN fragment;
-- +goose StatementEnd
