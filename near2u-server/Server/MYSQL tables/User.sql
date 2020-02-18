create table User (
name varchar(30) NOT NULL, 
surname varchar(30) NOT NULL, 
email varchar(50) NOT NULL UNIQUE, 
password varchar(30) NOT NULL, 
auth_token varchar(30), 
ID int NOT NULL PRIMARY KEY AUTO_INCREMENT);