create table Sensore (
name varchar(30) NOT NULL,  
type varchar(30) NOT NULL,
cod_sensore int NOT NULL PRIMARY KEY);


    delimiter //
        create procedure Sensore_insert (ambiente_cod varchar(40), sensor_name varchar(30), sensor_cod int, sensor_type varchar(30))
        begin
        
        insert into Sensore (name,type,cod_sensore) values (sensor_name,sensor_type,sensor_cod);
        insert into Sensore_Ambiente (cod_ambiente,cod_sensore) values (ambiente_cod,sensor_cod);
        end//

        delimiter ;