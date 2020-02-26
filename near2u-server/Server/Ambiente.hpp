#ifndef SERVER_AMBIENTE_H
#define SERVER_AMBIENTE_H

#include <string>
#include <iostream>
#include <assert.h>

#include "Sensore.h"

class Ambiente
{
	public:
	Ambiente( std::string name, std::string cod);
	std::string getNome();
	std::string getcodAmbiente();
	std::list<Sensore> * getSensori();
private:
	std::string Nome;

	std::string codAmbiente;

	std::list<Sensore> sensori;

};

#endif
