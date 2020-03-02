create table Dispositivo_Ambiente(
    cod_ambiente varchar(40) NOT NULL,
    code int NOT NULL,
    FOREIGN KEY (cod_ambiente) References Ambiente (cod_ambiente) on delete cascade on update cascade,
    FOREIGN KEY (code) References Dispositivo (code) on delete cascade on update cascade,
    primary key (cod_ambiente,code)
);