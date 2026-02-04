create table request_log
(
    id          bigserial   primary key,
    api_key     varchar     not null,
    ip          varchar     not null,
    endpoint    varchar     not null,
    timestamp   timestamptz not null default now(),
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);
