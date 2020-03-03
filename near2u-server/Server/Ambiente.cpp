

#include "Ambiente.hpp"


    Ambiente::Ambiente(std::string& name, std::string& cod_ambiente){
        Ambiente::Nome = name;
        Ambiente::codAmbiente = cod_ambiente;
    }
    std::string& Ambiente::getcodAmbiente(){
        return Ambiente::codAmbiente;
    }
    std::string& Ambiente::getNome(){
        return Ambiente::Nome;
    }
    std::list<Dispositivo> * Ambiente::getDispositivi(){
        return &dispositivi;
    }
    void Ambiente::addSensore(int code, std::string& nome, std::string& tipo){

        Sensore sensore(code,nome,tipo);
        dispositivi.push_back(sensore);

    }
    void Ambiente::addAttuatore(int code, std::string& nome, std::string& tipo){
        //Attuatore attuatore(code,nome,tipo);
        //dispositivi.push_back(attuatore);
    }

    Attuatore * Ambiente::getAttuatore(int cod_attuatore){
        std::list<Dispositivo>::iterator dispositivi_iterator;

        for(dispositivi_iterator=dispositivi.begin(); dispositivi_iterator != dispositivi.end(); dispositivi_iterator ++){
            if((*dispositivi_iterator).getCodice() == cod_attuatore){
                if(typeid(dispositivi_iterator) == typeid(Sensore)){
                    return static_cast<Attuatore*>(&(*dispositivi_iterator));
                }
            }
        }
        return nullptr;

    }
    void Ambiente::addComando(int code, std::string& comando){
        Attuatore * attuatore = Ambiente::getAttuatore(code);
        if(attuatore != nullptr){
            attuatore->getComandi()->push_back(comando);
        }

    }
    Sensore * Ambiente::getSensore(int cod_sensore){
        std::list<Dispositivo>::iterator dispositivi_iterator;

        for(dispositivi_iterator=dispositivi.begin(); dispositivi_iterator != dispositivi.end(); dispositivi_iterator ++){
            if((*dispositivi_iterator).getCodice() == cod_sensore){
                if(typeid(dispositivi_iterator) == typeid(Sensore)){
                    return static_cast<Sensore*>(&(*dispositivi_iterator));
                }
            }
        }
        return nullptr;


    }


    void Ambiente::deleteDispositivo(int cod_dispositivo){

        std::list<Dispositivo>::iterator dispositivi_iterator;

        for(dispositivi_iterator=dispositivi.begin(); dispositivi_iterator != dispositivi.end(); dispositivi_iterator ++){
            if((*dispositivi_iterator).getCodice() == cod_dispositivo){
                dispositivi.erase(dispositivi_iterator);
                break;
            }
        }
        
    }

    void Ambiente::addDispositivo(int code, std::string& nome, std::string& tipo, std::list<std::string> * commands){

        if(commands != nullptr){
            Attuatore dispositivo(code,nome,tipo,commands);
            dispositivi.push_back(dispositivo);
        }
        else {
            Sensore dispositivo(code,nome,tipo);
            dispositivi.push_back(dispositivo);
        }
       
        
    }


  
