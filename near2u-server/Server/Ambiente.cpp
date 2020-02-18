

#include "Ambiente.h"

namespace Server
{
    Ambiente::Ambiente(std::string name, int cod_ambiente){
        Ambiente::Nome = name;
        Ambiente::codAmbiente = cod_ambiente;
    }
    int Ambiente::getcodAmbiente(){
        return Ambiente::codAmbiente;
    }
    std::string Ambiente::getNome(){
        return Ambiente::Nome;
    }
    std::list<Sensore> * Ambiente::getSensori(){
        return &sensori;
    }


}  
