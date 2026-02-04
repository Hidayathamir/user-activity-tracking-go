alter table client
add constraint unique_client_api_key UNIQUE (api_key);
