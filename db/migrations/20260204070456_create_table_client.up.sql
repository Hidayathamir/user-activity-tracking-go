create table client
(
    id          bigserial   primary key,
    name        varchar     not null,
    email       varchar     not null,
    api_key     varchar     not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);
