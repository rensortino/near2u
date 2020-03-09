#ifndef SERVER_AMBIENTE_H
#define SERVER_AMBIENTE_H

#include <string>
#include <iostream>
#include <assert.h>

#include "Sensore.h"
#include "Attuatore.hpp"
#include "MQTTClient.h"

#define ADDRESS     "tcp://localhost:8082"
#define CLIENTSERVER    "server"
#define QOS         1
#define TIMEOUT     10000L

class Ambiente
{
	public:
<<<<<<< HEAD
	Ambiente( std::string name, std::string cod);
	std::string getNome();
	std::string getcodAmbiente();
	std::list<Sensore> * getSensori();
	Sensore * getSensore(int cod_sensore);
=======
	Ambiente( std::string& name, std::string& cod);
	~Ambiente();
	std::string& getNome();
	std::string& getcodAmbiente();
	std::list<Dispositivo *> *  getDispositivi();
	void deleteDispositivo(int cod_dispositivo);
	void addDispositivo(int code, std::string& nome, std::string& tipo, std::list<std::string> * commands);
	Dispositivo * getDispositivo(int code);
	bool inviaComando(int code_attuatore, std::string& comando);
>>>>>>> Iterazione_3
private:

	std::string Nome;

	std::string codAmbiente;

	std::list< Dispositivo *> dispositivi;

};

#endif
