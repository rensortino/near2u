create table Attuatore (

code int,

FOREIGN KEY (code) References Dispositivo (code) on delete cascade on update cascade,
primary key (code)
);


  delimiter //
        create procedure Attuatore_insert (ambiente_cod varchar(40), attuator_name varchar(30), attuator_cod int, attuator_type varchar(30))
        begin

        insert into Dispositivo (name,type,code) values (attuator_name,attuator_type,attuator_cod);
        insert into Attuatore (cod_attuatore) values (attuator_cod);
        insert into Dispositivo_Ambiente (cod_ambiente,code) values (ambiente_cod,attuator_cod);
        end//

        delimiter ;