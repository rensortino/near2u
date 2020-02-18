create table Ambiente_User(
    cod_ambiente int NOT NULL,
    User_id int NOT NULL,
    FOREIGN KEY (cod_ambiente) References Ambiente (cod_ambiente) on delete cascade on update cascade,
    FOREIGN KEY (User_id) References User (ID) on delete cascade on update cascade,
    primary key (cod_ambiente,User_id)
);

insert into Ambiente_User (cod_ambiente,User_id) values (int,int);

insert into Ambiente (name,cod_ambiente) values ('string',int);