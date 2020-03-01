

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
    std::list<Sensore> * Ambiente::getSensori(){
        return &sensori;
    }
    void Ambiente::addSensore(int code, std::string& nome, std::string& tipo){

        Sensore sensore(code,nome,tipo);
        sensori.push_back(sensore);

    }
    Sensore * Ambiente::getSensore(int cod_sensore){
        std::list<Sensore>::iterator sensori_iterator;

        for(sensori_iterator=sensori.begin(); sensori_iterator != sensori.end(); sensori_iterator ++){
            if((*sensori_iterator).getCodSensore() == cod_sensore){
                return &(*sensori_iterator);
            }
        }
        return nullptr;


    }


    void Ambiente::deleteSensore(int cod_sensore){

        std::list<Sensore>::iterator sensori_iterator;

        for(sensori_iterator=sensori.begin(); sensori_iterator != sensori.end(); sensori_iterator ++){
            if((*sensori_iterator).getCodSensore() == cod_sensore){
                sensori.erase(sensori_iterator);
                break;
            }
        }
        
    }


  
