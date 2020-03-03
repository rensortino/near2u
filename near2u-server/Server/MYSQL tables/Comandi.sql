create table Comandi (

comando varchar(30) NOT NULL,
cod_attuatore int,
FOREIGN KEY (cod_attuatore) References Attuatore (code) on delete cascade on update cascade,
primary key(comando,cod_attuatore)
);

insert into Comandi (comando,cod_attuatore) values ('',);