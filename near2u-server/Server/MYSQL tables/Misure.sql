create table Misure (
misura float NOT NULL,
code int,
ID int NOT NULL PRIMARY KEY AUTO_INCREMENT,
FOREIGN KEY (code) References Sensore (code) on delete cascade on update cascade
);

insert into Misure (misura,code) values (misura,code);