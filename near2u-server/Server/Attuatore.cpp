#include "Attuatore.hpp"

Attuatore::Attuatore(int code, std::string& type, std::string& name,std::list<std::string> * commands) : Dispositivo(type,name,code){
     comandi = *commands;

}

std::list<std::string> * Attuatore::getComandi(){
    return &comandi;
}


device_type Attuatore::get_device_type(){
    return device_type::attuatore;
}

bool Attuatore::controllaComando(std::string& comando){

    for(std::string comand : comandi){
        if(comand.compare(comando) == 0){
            return true;
        }
    }
    return false;

}