#ifndef CONTROLLER_H
#define CONTROLLER_H

#include <string>
#include <iostream>
#include <assert.h>
#include <jsoncpp/json/json.h>
#include "Ambiente.hpp"
#include "User.hpp"
#include <shared_mutex>
#include "function_mqtt.hpp"


class Controller
{
private:
	std::list<User *> users;
	static Controller* instance;
	std::shared_mutex User_mutex;
	MQTTClient client;

	User* Auth(std::string auth_token);
	User * search_on_cache(std::string email,std::string password);


public:
	void setUpMqtt();
	std::list<User *> *  getUsers();
	std::shared_mutex * getUser_mutex();
	static Controller* getIstance();
	Json::Value Topic_Ambiente(Json::Value data);
	Json::Value Register(Json::Value data);
	Json::Value Login(Json::Value data);
	Json::Value Crea_Ambiente(Json::Value data);
	Json::Value Associa_Utente(Json::Value data);
	Json::Value Visualizza_Ambienti(Json::Value data);
	Json::Value Inserisci_Dispositivi(Json::Value data);
	Json::Value Visualizza_Dispositivi(Json::Value data);
	Json::Value Elimina_Dispositivi(Json::Value data);
	Json::Value Invia_Comando(Json::Value data);
	Json::Value Visualizza_Storico(Json::Value data);
	Json::Value Elimina_Ambiente(Json::Value data);
	Json::Value Logout(Json::Value data);
	~Controller();


	

};
#endif
