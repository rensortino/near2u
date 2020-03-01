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

	std::list<User> * Controller::getUsers(){
		return &users;
	}
	std::shared_mutex * Controller::getUser_mutex(){
		return &User_mutex;
	}
	
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
		if(MYSQL::Query(query) == 0){
			response["status"] = "Succesfull";
			response["data"]["name"] = data["name"].asString();
			response["data"]["surname"] = data["surname"].asString();
			response["data"]["email"] = data["email"].asString();
			response["error"] = "";
		}
		else {
			response["status"] = "Failed";
			response["data"]["name"] = "";
			response["error"] = "Error in registration check the credential or contact the system admin";
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
					user.setAdmin(true);// to give admin privilege
					std::string query = "select Ambiente.cod_ambiente, Ambiente.name from (Ambiente join Ambiente_User on Ambiente.cod_ambiente = Ambiente_User.cod_ambiente) where User_email = '"+ user.getemail()+"'; ";
					sql::ResultSet *ambienti_db = MYSQL::Select_Query(query);
					if(ambienti_db ->rowsCount() > 0){
						while(ambienti_db->next()){
							std::string name = ambienti_db->getString("name");
							std::string cod_ambiente = ambienti_db->getString("cod_ambiente");
							Ambiente ambiente(name ,cod_ambiente);
							std::string query = "select Sensore.name, Sensore.type, Sensore.cod_sensore from (Sensore join Sensore_Ambiente on Sensore.cod_sensore = Sensore_Ambiente.cod_sensore) where Sensore_Ambiente.cod_ambiente = '"+ambiente.getcodAmbiente() + "';";
							sql::ResultSet *sensori_db = MYSQL::Select_Query(query);
							if (sensori_db->rowsCount() > 0){
								while(sensori_db->next()){
									std::string nome=sensori_db->getString("name");
									std::string tipo= sensori_db ->getString("type");
									int code= sensori_db->getInt("cod_sensore");
									Sensore sensore(code, nome, tipo);
									ambiente.getSensori()->push_back(sensore);
								}
							}
							user.getAmbienti()->push_back(ambiente);
						}
					}

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
	Json::Value Controller::Topic_Ambiente(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());

		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
			User_mutex.lock_shared();
			std::string cod_ambiente = Current_User ->getemail() + data["data"]["name"].asString();
			Ambiente * ambiente = Current_User->getAmbiente(cod_ambiente);
			if(ambiente != nullptr){
				response["status"] = "Succesfull";
				response["data"]["broker_host"] = "localhost:8082"; // qua poi inserire una variabile d'ambiente
				response["data"]["topic"] =  ambiente->getcodAmbiente();
			}
			else
			{
				response["status"] = "Failed";
				response["error"] = "Ambiente Not Found";
				response["data"] = "";
			}
			User_mutex.unlock_shared();
			return response;
		
		
		}

	
	
	
	
	
	
	
	Json::Value Controller::Configura_ambiente(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string cod_ambiente = Current_User ->getemail() + data["data"]["name"].asString();
		std::string name = data["data"]["name"].asString();

		
		std::string query = "call Ambiente_insert ('"+name +"','"+ Current_User->getemail() + "','"+ cod_ambiente+"');";
		
		if (MYSQL::Query(query) == 0){
				response["data"]["name"] = name;
				response["status"] = "Succesfull";
				response["error"] = "";
				Controller::User_mutex.lock();
				Current_User->addAmbiente(name,cod_ambiente);
				Controller::User_mutex.unlock();
		}
			else {
				response["data"] = "";
				response["status"] = "Failed";
				response["error"] = "Error in creating new Ambiente";
			}

		return response;
	}

	Json::Value Controller::Inserisci_Sensori(Json::Value data){
		Json::Value response;
		std::list<std::string> transaction;
		User * Current_User = Controller::Auth(data["auth"].asString());
		bool succes_flag = true;

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string cod_ambiente = Current_User ->getemail() + data["data"]["envname"].asString();
		 auto entriesArray = data["data"]["sensors"];
		 Json::Value::iterator sensors_to_add;
		 std::string start_transaction = "START TRANSACTION;";
		 transaction.push_back(start_transaction);
		 for (sensors_to_add = entriesArray.begin(); sensors_to_add != entriesArray.end();sensors_to_add ++){
			
			std::string query = "insert into Sensore (name,type,cod_sensore) values ('"+(*sensors_to_add)["name"].asString() +"','"+(*sensors_to_add)["kind"].asString() +"',"+std::to_string((*sensors_to_add)["code"].asInt()) + ");";
			std::string query_1 = "insert into Sensore_Ambiente (cod_ambiente,cod_sensore) values ('"+ cod_ambiente +"',"+ std::to_string((*sensors_to_add)["code"].asInt())+");";
			transaction.push_back(query);
			transaction.push_back(query_1);
		 }
		 	std::string commit = "commit;";
		 	transaction.push_back(commit);
			
			if (MYSQL::Queries(transaction) ==  true){
				for (sensors_to_add = entriesArray.begin(); sensors_to_add != entriesArray.end();sensors_to_add ++){
					std::string nome = (*sensors_to_add)["name"].asString();
					int code = (*sensors_to_add)["code"].asInt();
					std::string tipo = (*sensors_to_add)["kind"].asString();
					Controller::User_mutex.lock();
					Current_User->addSensore(cod_ambiente,code,nome,tipo);
					Controller::User_mutex.unlock();
				}
			}
			else {
				succes_flag = false;
				response["status"] = "Failed";
				response["error"] = "Error in creating sensors ";
				response["data"] = "";
			}
			 
		
		if (succes_flag == true){
			response["status"] = "Succesfull";
			response["error"] = "";
			response["data"] = "insert completed";
		}
		 

		 return response;


	}

	Json::Value Controller::Visualizza_Ambienti(Json::Value data){
		Json::Value response;

		User * Current_User = Controller::Auth(data["auth"].asString());

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::list<Ambiente> * ambienti = Current_User->getAmbienti(); 
		if(ambienti->empty()){
			response["status"] = "Failed";
			response["error"] = "User does not have Ambiente associated";
			response["data"] = "";
			return response;
		}

		std::list<Ambiente>::iterator sensori_iterator;
			
		int i = 0;
		for(sensori_iterator = ambienti->begin();sensori_iterator != ambienti->end();sensori_iterator ++) {	
				
			response["data"]["environments"][i] = (*sensori_iterator).getNome();
			i++;
    	}
        response["status"] = "Succesfull";
		response["error"]= "";
		return response;
				
	}
			

	Json::Value Controller::Visualizza_Sensori(Json::Value data){
		Json::Value response;

		User * Current_User = Controller::Auth(data["auth"].asString());
		std::string cod_ambiente = Current_User ->getemail() + data["data"]["envname"].asString();
		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		int i = 0;
		std::list<Sensore> * sensori = Current_User->getSensori(cod_ambiente);
		if(sensori == nullptr){
			response["status"] = "Failed";
			response["error"] = "Ambiente Not Found";
			response["data"] = "";
			return response;
		}
		response["status"] = "Succesfull";
		response["data"]["sensors"]= Json::Value(Json::arrayValue);
		std::list<Sensore>::iterator sensori_iterator;
		for(sensori_iterator = sensori->begin();sensori_iterator != sensori->end();sensori_iterator ++)
		{	
			Json::Value sensor;
			sensor["name"] =(*sensori_iterator).getName();
			sensor["kind"] = (*sensori_iterator).getType();
			sensor["code"] = (*sensori_iterator).getCodSensore();
			response["data"]["sensors"][i] = sensor;
			i++;
		}
        response["status"] = "Succesfull";
				
		return response;
		

	}

	Json::Value Controller::Elimina_sensori(Json::Value data){
		Json::Value response;
		std::list<std::string> transaction;
		User * Current_User = Controller::Auth(data["auth"].asString());

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string cod_ambiente = Current_User ->getemail() + data["data"]["envname"].asString();
		auto entriesArray = data["data"]["sensors"];
		Json::Value::iterator sensors_to_delete;
		std::string start_transaction = "START TRANSACTION;";
		transaction.push_back(start_transaction);
		for (sensors_to_delete = entriesArray.begin(); sensors_to_delete != entriesArray.end();sensors_to_delete ++){
			
			std::string query = "delete from Sensore where cod_sensore = "+std::to_string((*sensors_to_delete).asInt()) + ";";
			transaction.push_back(query);
		}
		std::string commit = "commit;";
		transaction.push_back(commit);
			
		if (MYSQL::Queries(transaction) ==  true){
			for (sensors_to_delete = entriesArray.begin(); sensors_to_delete != entriesArray.end();sensors_to_delete ++){
				Controller::User_mutex.lock();
				int cod_sensore = (*sensors_to_delete).asInt();
				Current_User->deleteSensore(cod_ambiente,cod_sensore);
				Controller::User_mutex.unlock();
			}
					
					
		}
		else{
			response["status"] = "Failed";
			response["error"] = "Error in deleting sensors";
			response["data"] = "";
			return response;
		}
		response["status"] = "Succesfull";
		response["error"] = "";
		response["data"] = "deletion completed";
		return response;
	}






		

	

