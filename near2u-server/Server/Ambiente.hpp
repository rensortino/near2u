#ifndef SERVER_AMBIENTE_H
#define SERVER_AMBIENTE_H

#include <string>
#include <iostream>
#include <assert.h>

#include "Sensore.h"
#include "Attuatore.hpp"

class Ambiente
{
	public:
	Ambiente( std::string& name, std::string& cod);
	std::string& getNome();
	std::string& getcodAmbiente();
	std::list<Dispositivo *> *  getDispositivi();
	void addSensore(int code, std::string& nome, std::string& tipo);
	void deleteDispositivo(int cod_dispositivo);
	void addDispositivo(int code, std::string& nome, std::string& tipo, std::list<std::string> * commands);

private:

	std::string Nome;

	std::string codAmbiente;

	std::list< Dispositivo *> dispositivi;

};

#endif
