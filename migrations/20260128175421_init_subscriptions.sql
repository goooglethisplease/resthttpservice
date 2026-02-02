-- +goose Up
-- +goose StatementBegin
create extension if not exists "uuid-ossp";

create table if not exists subscriptions (
    id uuid primary key default uuid_generate_v4(),
    service_name text not null,
    price int not null check (price > 0),
    user_id uuid not null,

    start_date date not null,
    end_date date null,

    check (end_date is null or end_date >= start_date)
);

create index if not exists idx_subscriptions_user_id
    on subscriptions(user_id);

create index if not exists idx_subscriptions_user_service
    on subscriptions(user_id, service_name);

create index if not exists idx_subscriptions_start_date
    on subscriptions(start_date);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists subscriptions;
-- +goose StatementEnd
