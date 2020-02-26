#ifndef CONTROLLER_H
#define CONTROLLER_H

#include <string>
#include <iostream>
#include <assert.h>
#include <jsoncpp/json/json.h>

#include "Ambiente.hpp"
#include "User.hpp"
#include <shared_mutex>


class Controller
{
private:
	std::list<User> users;
	static Controller* instance;
	std::shared_mutex User_mutex;


	User* Auth(std::string auth_token);
	User * search_on_cache(std::string email,std::string password);


public:
	static Controller* getIstance();
	Json::Value Seleziona_Ambiente(Json::Value data);
	Json::Value Register(Json::Value data);
	Json::Value Login(Json::Value data);
	Json::Value Configura_ambiente(Json::Value data);
	Json::Value Inserisci_Sensori(Json::Value data);
	

};
#endif
