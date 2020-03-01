#include <string>
#include <vector>
#include <list>
#include <iostream>
#include <assert.h>

#include "Sensore.h"


Sensore::Sensore(int codice, std::string& nome,std::string& tipo){
    codSensore = codice;
    name = nome;
    type = tipo;
}

int Sensore::getCodSensore(){
    return codSensore;
}
std::string& Sensore::getName(){
    return name;
}
std::string& Sensore::getType(){
    return type;
}