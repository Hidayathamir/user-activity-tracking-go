alter table client_request_count
add constraint fk_client_request_count_api_key
foreign key (api_key)
references client(api_key)
on delete cascade;
