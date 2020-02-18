#ifndef CONTROLLER_H
#define CONTROLLER_H

#include <string>
#include <iostream>
#include <assert.h>
#include <jsoncpp/json/json.h>

#include "Ambiente.cpp"
#include "User.h"


namespace Server
{
class Controller
{
private:
	std::list<User> users;
	static Controller* instance;

	User* Auth(std::string auth_token);


public:
	static Controller* getIstance();
	Json::Value Seleziona_Ambiente(Json::Value data);
	Json::Value Register(Json::Value data);
	Json::Value Login(Json::Value data);

};
Controller *Controller::instance = 0;
}  // namespace Server
#endif
