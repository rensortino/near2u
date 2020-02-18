#ifndef SERVER_AMBIENTE_H
#define SERVER_AMBIENTE_H

#include <string>
#include <iostream>
#include <assert.h>

#include "./after/Sensore.h"

namespace Server
{
class Ambiente
{
	public:
	Ambiente( std::string name, int cod);
	std::string getNome();
	int getcodAmbiente();
	std::list<Sensore> * getSensori();
private:
	std::string Nome;

	int codAmbiente;

	std::list<Sensore> sensori;

};

}  
#endif
