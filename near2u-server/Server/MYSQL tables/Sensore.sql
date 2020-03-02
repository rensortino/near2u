create table Sensore (

code int,
FOREIGN KEY (code) References Dispositivo (code) on delete cascade on update cascade,
primary key (code)
);



    delimiter //
        create procedure Sensore_insert (ambiente_cod varchar(40), name varchar(30), cod int, type varchar(30))
        begin

        insert into Dispositivo (name,type,code) values (name,type,cod);
        insert into Sensore (cod_sensore) values (sensor_cod);
        insert into Dispositivo_Ambiente (cod_ambiente,code) values (ambiente_cod,sensor_cod);
        end//

        delimiter ;