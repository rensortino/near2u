CREATE TABLE Ambiente (
name varchar(30) NOT NULL,  
cod_ambiente varchar(40) not null  ,
primary key (cod_ambiente)
);



    delimiter //
    create procedure Ambiente_insert (Ambiente_name varchar(30), email varchar(50), cod varchar(40))
    begin
    
    insert into Ambiente (name,cod_ambiente) values (Ambiente_name,cod);
    insert into Ambiente_User (cod_ambiente,User_email) values (cod,email);
    end//

    delimiter ;