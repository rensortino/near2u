#include <string>
#include <iostream>
#include <assert.h>

#include "Controller.hpp"
#include "SHA_CRYPTO.hpp"
#include "MYSQL.hpp"

	



	Controller* Controller::getIstance(){
		if (!instance)
      		instance = new Controller;
      return instance;

	}
	Controller *Controller::instance = 0;

	User * Controller::search_on_cache(std::string email,std::string password){
		std::list<User>::iterator cache_user;
		User_mutex.lock_shared();
			for (cache_user = users.begin(); cache_user != users.end(); ++cache_user){
				if((cache_user->getemail().compare(email) == 0) && (cache_user->getPassword().compare(password) == 0) ){
					User_mutex.unlock_shared();
					return  &(*cache_user);
				}
			}
				User_mutex.unlock_shared();
		return nullptr;
	}
	Json::Value Controller::Register(Json::Value data){
		Json::Value response;
		

		std::cout << data.toStyledString() <<std::endl;
		std::string query = "INSERT INTO User (name, surname, email, password,auth_token)\
                                VALUES ('"+ data["name"].asString() + "','" + data["surname"].asString() + "','"+data["email"].asString()+"','"+data["password"].asString()+"','"+SHA_Crypto(data["email"].asString() + data["password"].asString())+"');";
		response = MYSQL::insert(query);
		if(response["status"].asString().compare("Succesfull") == 0){
			response["data"]["name"] = data["name"].asString();
			response["data"]["surname"] = data["surname"].asString();
			response["data"]["email"] = data["email"].asString();
		}
		  return response;
	
	}
	Json::Value Controller::Login(Json::Value data){
		
		Json::Value response;
		response["status"] = "";
		response["error"] = "";
		
		
		
		if(search_on_cache(data["Email"].asString(),data["Password"].asString()) == nullptr){
			std::string query = "select name,surname,email,auth_token,password,Admin from User where email = '" + data["Email"].asString() + "' and password = '" + data["Password"].asString() + "'";
			sql::ResultSet  *res;
			res = MYSQL::Select_Query(query);
			if( res->rowsCount() == 0){
				response["status"] = "failed";
				response["error"] = "No user found please check credentials";
				response["data"] = "";
			}
			else {
				while (res->next()) {
					std::cout << "auth_token = '" << res->getString("auth_token") << "'" << std::endl;
					response["data"]["auth"] = (std::string) res->getString("auth_token");
					User user((std::string) res->getString("name"),(std::string) res->getString("surname"),(std::string) res->getString("email"),(std::string) res->getString("auth_token"),(std::string) res->getString("password"));
					user.setAdmin(res->getBoolean("Admin"));
					User_mutex.lock();
					Controller::users.push_back(user);
					User_mutex.unlock();
				}
				response["status"] = "Succesfull";   
			}
			delete res;
		}
		else{
			response["status"] = "Succesfull";
			response["data"]["auth"] = search_on_cache(data["Email"].asString(),data["Password"].asString())->getauth_token();
		}
			return response;  


	}
	User* Controller::Auth(std::string auth_token){
		std::list<User>::iterator cache_user;
		User_mutex.lock_shared();
		for (cache_user = users.begin(); cache_user != users.end(); ++cache_user){
    		if(cache_user->getauth_token().compare(auth_token) == 0){
				User_mutex.unlock_shared();
				return  &(*cache_user);
			}
		}
		User_mutex.unlock_shared();

		return nullptr;
		
		


		
	}
	Json::Value Controller::Seleziona_Ambiente(Json::Value data){
		Json::Value response;
		response["status"] = "";
		response["error"] = "";
		User * Current_User = Controller::Auth(data["auth"].asString());

		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
	

		std::list<Ambiente>::iterator cache_ambiente;
			User_mutex.lock_shared();
			for (cache_ambiente = Current_User->getAmbienti()->begin(); cache_ambiente != Current_User->getAmbienti()->end(); ++cache_ambiente){
				if(cache_ambiente->getcodAmbiente().compare(Current_User ->getemail() + data["data"]["name"].asString()) == 0 ){
					response["status"] = "Succesfull";
					response["data"]["broker_host"] = "localhost:8082"; // qua poi inserire una variabile d'ambiente
					response["data"]["topic"] =  cache_ambiente->getcodAmbiente();
					User_mutex.unlock_shared();
					return response;
				}
			}
			User_mutex.unlock_shared();

			std::string query = "select Ambiente.cod_ambiente, Ambiente.name from User join (Ambiente_User join Ambiente on Ambiente.cod_ambiente = Ambiente_User.cod_ambiente ) on User_email = User.email  where Ambiente.name = '"+ data["data"]["name"].asString() + "' and User.email ='"+Current_User->getemail() +"';"; 
			std::cout << query << std::endl;
			sql::ResultSet *res = MYSQL::Select_Query(query);

			if( res == nullptr || res->rowsCount() == 0 ){
            response["status"] = "Failed";
            response["error"] = "Ambiente not Found";
        	}
			else
			{
				while (res->next()) {
					response["data"]["topic"] = (std::string) res->getString("cod_ambiente");
					Ambiente ambiente((std::string) res->getString("name"), res->getString("cod_ambiente"));
					User_mutex.lock();
					Current_User->getAmbienti()->push_back(ambiente);
					User_mutex.unlock();
            }
            	response["status"] = "Succesfull";
				response["data"]["broker_host"] = "localhost:8082"; 
				
			}
			delete res;
			return response;
		
		
		}

	
	
	
	
	
	
	
	Json::Value Controller::Configura_ambiente(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		std::string cod_ambiente = Current_User ->getemail() + data["data"]["name"].asString();

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}

		
		std::string query = "call Ambiente_insert ('"+data["data"]["name"].asString() +"','"+ Current_User->getemail() + "','"+ cod_ambiente+"');";
		response = MYSQL::insert(query);
		

		
			if (response["status"].asString().compare("Succesfull") == 0){
				response["data"]["name"] = data["data"]["name"].asString();
				Ambiente ambiente(data["data"]["name"].asString(),cod_ambiente);
				Controller::User_mutex.lock();
				Current_User->getAmbienti()->push_back(ambiente);
				Controller::User_mutex.unlock();
			}
			
		
		


		return response;
	}

	Json::Value Controller::Inserisci_Sensori(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		std::string cod_ambiente = Current_User ->getemail() + data["data"]["envName"].asString();

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		 auto entriesArray = data["data"]["sensors"];
		 Json::Value::iterator sensors_to_add;
		 for (sensors_to_add = entriesArray.begin(); sensors_to_add != entriesArray.end();sensors_to_add ++){
			
			std::string query = "call Sensore_insert ('"+cod_ambiente +"','"+ (*sensors_to_add)["name"].asString() + "'," + std::to_string((*sensors_to_add)["code"].asInt()) +",'"+ (*sensors_to_add)["kind"].asString()+ "');";
			response = MYSQL::insert(query);

			if (response["status"].asString().compare("Succesfull") == 0){
				Sensore sensore((*sensors_to_add)["code"].asInt(),(*sensors_to_add)["name"].asString(),(*sensors_to_add)["kind"].asString());
				Controller::User_mutex.lock();
				Ambiente * ambiente = Current_User->getAmbiente(cod_ambiente);
				ambiente->getSensori()->push_back(sensore);
				Controller::User_mutex.unlock();
			}
			else {
				break;
			}
			 
			 
		 }

		 return response;


	}





		

	

