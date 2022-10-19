-- Write your migrate up statements here
create table status_action
(
    id   serial primary key,
    name varchar(40) NOT NULL

);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
