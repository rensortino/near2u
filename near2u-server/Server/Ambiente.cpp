

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
    std::list<Dispositivo *> * Ambiente::getDispositivi(){
        return &dispositivi;
    }
    void Ambiente::addSensore(int code, std::string& nome, std::string& tipo){

        Sensore sensore(code,nome,tipo);
        dispositivi.push_back(&sensore);

    }


    void Ambiente::deleteDispositivo(int cod_dispositivo){

        std::list<Dispositivo *>::iterator dispositivi_iterator;

        for(dispositivi_iterator=dispositivi.begin(); dispositivi_iterator != dispositivi.end(); dispositivi_iterator ++){
            
            if((*dispositivi_iterator)->getCodice() == cod_dispositivo){
                delete (*dispositivi_iterator);
                dispositivi.erase(dispositivi_iterator);
                break;
            }
        }
        
    }

    void Ambiente::addDispositivo(int code, std::string& nome, std::string& tipo, std::list<std::string> * commands){

        if(commands != nullptr){
            Attuatore *  dispositivo = new Attuatore(code,nome,tipo,commands);
            dispositivi.push_back(dispositivo);
        }
        else {
            Sensore *  dispositivo = new Sensore(code,nome,tipo);
            dispositivi.push_back(dispositivo);
        }
       
        
    }

    Dispositivo * Ambiente::getDispositivo(int code){

        std::list<Dispositivo *>::iterator dispositivi_iterator;

        for(dispositivi_iterator=dispositivi.begin(); dispositivi_iterator != dispositivi.end(); dispositivi_iterator ++){
            
            if((*dispositivi_iterator)->getCodice() == code){
                return (*dispositivi_iterator);
                
            }
        }

        return nullptr;


    }

    bool Ambiente::inviaComando(int code_attuatore, std::string& comando){
        Attuatore * attuatore;
        try{
            attuatore = static_cast<Attuatore *>(getDispositivo(code_attuatore));
        }
        catch(std::exception &error){
            std::cout << error.what() << std::endl;
            return false;
        }

        if(attuatore->controllaComando(comando) == false){
            return false;
        }
        return true;
        


    }


  
