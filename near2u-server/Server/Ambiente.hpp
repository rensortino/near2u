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
	std::list<Dispositivo> * getDispositivi();
	void addSensore(int code, std::string& nome, std::string& tipo);
	Sensore * getSensore(int cod_sensore);
	Attuatore * getAttuatore(int cod_attuatore);
	void deleteDispositivo(int cod_dispositivo);
	void addAttuatore(int code, std::string& nome, std::string& tipo);
	void addComando(int code, std::string& comando);
	void addDispositivo(int code, std::string& nome, std::string& tipo, std::list<std::string> * commands);

private:

	std::string Nome;

	std::string codAmbiente;

	std::list<Dispositivo> dispositivi;

};

#endif
