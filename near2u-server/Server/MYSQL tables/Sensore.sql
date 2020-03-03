create table Sensore (

code int,
FOREIGN KEY (code) References Dispositivo (code) on delete cascade on update cascade,
primary key (code)
);




        insert into Dispositivo (name,type,code) values (name,type,cod);
        insert into Sensore (cod_sensore) values (sensor_cod);
        insert into Dispositivo_Ambiente (cod_ambiente,code) values (ambiente_cod,sensor_cod);
        
