#include <string>
#include <iostream>
#include <assert.h>

#include "Controller.h"
#include "MYSQL.cpp"
#include "SHA_CRYPTO.cpp"

namespace Server
{

	Controller* Controller::getIstance(){
		if (!instance)
      		instance = new Controller;
      return instance;

	}
	Json::Value Controller::Register(Json::Value data){

		std::cout << data.toStyledString() <<std::endl;
		std::string query = "INSERT INTO User (name, surname, email, password,auth_token)\
                                VALUES ('"+ data["name"].asString() + "','" + data["surname"].asString() + "','"+data["email"].asString()+"','"+data["password"].asString()+"','"+SHA_Crypto(data["email"].asString() + data["password"].asString())+"');";
		return  MYSQL::insert(query);
	
	}
	Json::Value Controller::Login(Json::Value data){

		sql::ResultSet  *res;
		Json::Value response;
		
		std::string query = "select name,surname,email,auth_token from User where email = '" + data["email"].asString() + "' and password = '" + data["password"].asString() + "'";
		std::cout << query << std::endl;

		res = MYSQL::Select_Query(query);
		if( res->rowsCount() == 0){
            response["Status"] = "Failed";
            response["error"] = "No user found please check credentials";
        }
		if( res->rowsCount() > 1){
            response["Status"] = "Failed";
            response["error"] = "Ambiguity in User registration";
        }
		else {
            while (res->next()) {
                std::cout << ", auth_token = '" << res->getString("auth_token") << "'" << std::endl;
                response["auth"] = (std::string) res->getString("auth_token");
				User user((std::string) res->getString("name"),(std::string) res->getString("surname"),(std::string) res->getString("email"),(std::string) res->getString("auth_token"));
				Controller::users.push_back(user);
            }
            response["Status"] = "Succesfull";   
        }
		delete res;
		return response;  


	}
	User* Controller::Auth(std::string auth_token){
		std::list<User>::iterator cache_user;
		for (cache_user = users.begin(); cache_user != users.end(); ++cache_user){
    		if(cache_user->getauth_token().compare(auth_token) == 0){
				return  &(*cache_user);
			}
		}

		return nullptr;
		
		


		
	}
	Json::Value Controller::Seleziona_Ambiente(Json::Value data){
		Json::Value response;
		std::cout << "seleziona_ambiente"<< std::endl;
		User * Current_User = Controller::Auth(data["auth_token"].asString());

		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["message"] = "Unauthorized";
			return response;
		}
	

		std::list<Ambiente>::iterator cache_ambiente;

			for (cache_ambiente = Current_User->getAmbienti()->begin(); cache_ambiente != Current_User->getAmbienti()->end(); ++cache_ambiente){
				if(cache_ambiente->getNome().compare(data["name"].asString()) == 0 ){
					response["broker_host"] = "localhost:8082"; // qua poi inserire una variabile d'ambiente
					response["topic"] = cache_ambiente->getcodAmbiente();
					return response;
				}
			}


			std::string query = "select Ambiente.cod_ambiente, Ambiente.name from User join (Ambiente_User join Ambiente on Ambiente.cod_ambiente = Ambiente_User.cod_ambiente ) on User_id = User.ID  where Ambiente.name = '"+ data["name"].asString() + "';"; 
			sql::ResultSet *res = MYSQL::Select_Query(query);

			if( res->rowsCount() == 0){
            response["Status"] = "Failed";
            response["error"] = "Ambiente not Found";
        	}
			else
			{
				while (res->next()) {
					response["topic"] = res->getInt("cod_ambiente");
					Ambiente ambiente((std::string) res->getString("name"), res->getInt("cod_ambiente"));
					Current_User->getAmbienti()->push_back(ambiente);
            }
            	response["Status"] = "Succesfull";
				response["broker_host"] = "localhost:8082"; 
				
			}
			delete res;
			return response;
			



		
		
		}




		
	}

	

