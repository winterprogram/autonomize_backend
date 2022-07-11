-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_file
(
    id bigserial,
    user_id integer,
    file_name text,
    url text,
    created_at timestamp without time zone DEFAULT NOW(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
