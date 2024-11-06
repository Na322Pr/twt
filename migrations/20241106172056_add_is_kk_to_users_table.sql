-- +goose Up
-- +goose StatementBegin
alter table users ADD COLUMN is_kk boolean;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop COLUMN is_kk;
-- +goose StatementEnd
