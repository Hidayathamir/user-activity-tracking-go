alter table request_log
add constraint fk_request_log_api_key
foreign key (api_key)
references client(api_key)
on delete cascade;
