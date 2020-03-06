#ifndef SERVER_AMBIENTE_H
#define SERVER_AMBIENTE_H

#include <string>
#include <iostream>
#include <assert.h>

#include "Sensore.h"
#include "Attuatore.hpp"
#include "function_mqtt.hpp"


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
	Dispositivo * getDispositivo(int code);
	bool inviaComando(int code_attuatore, std::string& comando);
private:

	std::string Nome;

	std::string codAmbiente;

	std::list< Dispositivo *> dispositivi;

};

#endif
