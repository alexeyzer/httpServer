CREATE TABLE ref(
    id bigserial not null primary key,
    adv_id int not null,
    ref varchar(200) not null,
    constraint fk_adv
                foreign key(adv_id)
                    references adv(ID)
);