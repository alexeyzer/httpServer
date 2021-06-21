CREATE TABLE adv(
    id bigserial not null primary key,
    name varchar(200) not null,
    description varchar(1000),
    price int,
    date_create date default now()
);