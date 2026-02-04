create table client_request_count
(
    id          bigserial   primary key,
    api_key     varchar     not null,
    datetime    timestamptz not null default now(),
    count       int         not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);
