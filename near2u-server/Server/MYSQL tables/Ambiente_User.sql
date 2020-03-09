CREATE TABLE Ambiente_User(
    cod_ambiente varchar(40) NOT NULL,
    User_email varchar(50) NOT NULL,
    FOREIGN KEY (cod_ambiente) References Ambiente (cod_ambiente) on delete cascade on update cascade,
    FOREIGN KEY (User_email) References User (email) on delete cascade on update cascade,
    primary key (cod_ambiente,User_email)
);

