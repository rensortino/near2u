#include "Dispositivo.hpp"




Dispositivo::Dispositivo(std::string& type, std::string& name, int code){
        tipo=type;
        nome = name;
        codice = code;
    
}

std::string& Dispositivo::getTipo(){
    return tipo;
}
std::string& Dispositivo::getNome(){
    return nome;
}

int Dispositivo::getCodice(){
    return codice;
}

Dispositivo::~Dispositivo(){

}