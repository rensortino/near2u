

#include "Ambiente.hpp"


    Ambiente::Ambiente(std::string name, std::string cod_ambiente){
        Ambiente::Nome = name;
        Ambiente::codAmbiente = cod_ambiente;
    }
    std::string Ambiente::getcodAmbiente(){
        return Ambiente::codAmbiente;
    }
    std::string Ambiente::getNome(){
        return Ambiente::Nome;
    }
    std::list<Sensore> * Ambiente::getSensori(){
        return &sensori;
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


  
