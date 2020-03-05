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

    std::list<std::string>::iterator comandi_iterator;

    for(comandi_iterator = comandi.begin(); comandi_iterator != comandi.end(); comandi_iterator ++){
        if((*comandi_iterator).compare(comando) == 0){
            return true;
        }
    }
    return false;

}