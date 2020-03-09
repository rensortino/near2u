CREATE TABLE Sensore (

code int,
FOREIGN KEY (code) References Dispositivo (code) on delete cascade on update cascade,
primary key (code)
);




        insert into Dispositivo (name,type,code,cod_ambiente) values (name,type,cod,cod_ambiente);
        insert into Sensore (cod_sensore) values (sensor_cod);
        
