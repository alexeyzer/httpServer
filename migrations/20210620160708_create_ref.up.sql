CREATE TABLE ref(
    ID bigserial not null primary key,
    advId int not null,
    ref varchar(200) not null,
    constraint fk_adv
                foreign key(advId)
                    references adv(ID)
);