-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id SERIAL NOT NULL PRIMARY KEY,
    redirect_path VARCHAR NOT NULL UNIQUE,
    scheme VARCHAR NOT NULL,
    host VARCHAR NOT NULL,
    path VARCHAR,
    query VARCHAR
);

CREATE INDEX redirect_path_idx ON urls USING HASH (redirect_path);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
