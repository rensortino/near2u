create table Sensore_Ambiente(
    cod_ambiente varchar(40) NOT NULL,
    cod_sensore int NOT NULL,
    FOREIGN KEY (cod_ambiente) References Ambiente (cod_ambiente) on delete cascade on update cascade,
    FOREIGN KEY (cod_sensore) References Sensore (cod_sensore) on delete cascade on update cascade,
    primary key (cod_ambiente,cod_sensore)
);