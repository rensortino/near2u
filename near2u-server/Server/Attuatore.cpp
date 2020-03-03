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