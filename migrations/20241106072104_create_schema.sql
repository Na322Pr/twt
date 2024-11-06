-- +goose Up
-- +goose StatementBegin
create table users (
    id bigint primary key,
    name varchar(100),
    surname varchar(100),
    seat int,
    status varchar(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "users" cascade;
-- +goose StatementEnd
