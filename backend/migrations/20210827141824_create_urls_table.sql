-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id SERIAL NOT NULL PRIMARY KEY,
    redirect_path VARCHAR NOT NULL UNIQUE,
    scheme VARCHAR NOT NULL,
    host VARCHAR NOT NULL,
    path VARCHAR,
    query VARCHAR
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
