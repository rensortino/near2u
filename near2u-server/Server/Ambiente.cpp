

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


  
