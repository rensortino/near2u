create table Attuatore (

code int,

FOREIGN KEY (code) References Dispositivo (code) on delete cascade on update cascade,
primary key (code)
);


  
        insert into Dispositivo (name,type,code,cod_ambiente) values (attuator_name,attuator_type,attuator_cod,cod_ambiente);
        insert into Attuatore (cod_attuatore) values (attuator_cod);
        