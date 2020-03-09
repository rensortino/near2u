CREATE TABLE Dispositivo (
name varchar(30) NOT NULL,  
type varchar(30) NOT NULL,
code int NOT NULL PRIMARY KEY,
cod_ambiente varchar(40) NOT NULL,
FOREIGN KEY (cod_ambiente) References Ambiente (cod_ambiente) on delete cascade on update cascade
);