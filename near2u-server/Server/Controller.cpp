#include <string>
#include <iostream>
#include <assert.h>

#include "Controller.hpp"
#include "SHA_CRYPTO.hpp"
#include "MYSQL.hpp"

	
#define ADDRESS     "tcp://localhost:8082"
#define QOS         1
#define TIMEOUT     10000L


	Controller* Controller::getIstance(){
		if (!instance){
			instance = new Controller;
		}
		return instance;
      		

	}

	void Controller::setUpMqtt(){
		using namespace MQTT;
		if(client == nullptr){
			MQTTClient_connectOptions conn_opts;
			conn_opts = MQTTClient_connectOptions_initializer;
			

			MQTTClient_create(&client, ADDRESS, "Server_sensorData",
			MQTTCLIENT_PERSISTENCE_NONE, NULL);
			conn_opts.keepAliveInterval = 20;
			conn_opts.cleansession = 1;

			int rc;

			if (( rc= MQTTClient_connect(client, &conn_opts)) != MQTTCLIENT_SUCCESS)
			{
				printf("Failed to connect, return code %d\n", rc);
				exit(-1);
			}

			MQTTClient_setCallbacks(client, NULL, connlost, UploadDataSensor,delivered);
			
		}
	}
	Controller *Controller::instance = 0;

	std::list<User *> * Controller::getUsers(){
		return &users;
	}
	std::shared_mutex * Controller::getUser_mutex(){
		return &User_mutex;
	}
	
	User * Controller::search_on_cache(std::string email,std::string password){
		std::list<User *>::iterator cache_user;
		User_mutex.lock_shared();
			for (cache_user = users.begin(); cache_user != users.end(); ++cache_user){
				if(((*cache_user)->getemail().compare(email) == 0) && ((*cache_user)->getPassword().compare(password) == 0) ){
					User_mutex.unlock_shared();
					return  (*cache_user);
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
			response["status"] = "Successful";
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
					response["data"]["admin"] = res->getBoolean("Admin");
					User * user = new User((std::string) res->getString("name"),(std::string) res->getString("surname"),(std::string) res->getString("email"),(std::string) res->getString("auth_token"),(std::string) res->getString("password"));
					user->setAdmin(res->getBoolean("Admin"));
					std::string query = "select Ambiente.cod_ambiente, Ambiente.name from (Ambiente join Ambiente_User on Ambiente.cod_ambiente = Ambiente_User.cod_ambiente) where User_email = '"+ user->getemail()+"'; ";
					sql::ResultSet *ambienti_db = MYSQL::Select_Query(query);
					if(ambienti_db ->rowsCount() > 0){
						while(ambienti_db->next()){
							std::string name_ambiente = ambienti_db->getString("name");
							std::string cod_ambiente = ambienti_db->getString("cod_ambiente");
							user->addAmbiente(name_ambiente,cod_ambiente);
							query = "select Dispositivo.name, Dispositivo.type, Dispositivo.code from (Dispositivo join Sensore on Sensore.code = Dispositivo.code) where Dispositivo.cod_ambiente = '"+ cod_ambiente +"';";
							sql::ResultSet *sensori_db = MYSQL::Select_Query(query);
							while(sensori_db->next()){
								std::string nome_sensore=sensori_db->getString("name");
								std::string tipo_sensore= sensori_db ->getString("type");
								int code_sensore= sensori_db->getInt("code");
								user->addDispositivo(cod_ambiente,code_sensore,nome_sensore,tipo_sensore,nullptr);
							}
							query = "select Dispositivo.name, Dispositivo.type, Dispositivo.code from (Dispositivo join Attuatore on Attuatore.code = Dispositivo.code) where Dispositivo.cod_ambiente = '"+ cod_ambiente +"';";
							sql::ResultSet *attuatori_db = MYSQL::Select_Query(query);
							while(attuatori_db->next()){
									std::list<std::string> comandi;
									std::string nome_attuatore=attuatori_db->getString("name");
									std::string tipo_attuatore= attuatori_db ->getString("type");
									int code_attuatore= attuatori_db->getInt("code");
									std::string query = " select comando from Comandi where cod_attuatore = "+ std::to_string(code_attuatore	) +";";
									sql::ResultSet *comandi_db = MYSQL::Select_Query(query);
									while(comandi_db->next()){
										std::string comando=comandi_db->getString("comando");
										comandi.push_back(comando);
									}
									user->addDispositivo(cod_ambiente,code_attuatore,nome_attuatore,tipo_attuatore,&comandi);
							}
							MQTTClient_subscribe(client, cod_ambiente.c_str(), QOS);
						}
					}

					User_mutex.lock();
					Controller::users.push_back(user);
					User_mutex.unlock();
				}
				response["status"] = "Successful";   
			}
			delete res;
		}
		else{
			User * current_user = search_on_cache(data["Email"].asString(),data["Password"].asString());
			response["status"] = "Successful";
			response["data"]["auth"] = current_user->getauth_token();
			response["data"]["admin"] = current_user->getAdmin();
			
		}
			return response;  


	}

	User* Controller::Auth(std::string auth_token){
		std::list<User *>::iterator cache_user;
		User_mutex.lock_shared();
		for (cache_user = users.begin(); cache_user != users.end(); ++cache_user){
    		if((*cache_user)->getauth_token().compare(auth_token) == 0){
				User_mutex.unlock_shared();
				return  (*cache_user);
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
			std::string cod_ambiente = data["data"]["envcode"].asString();
			Ambiente * ambiente = Current_User->getAmbiente(cod_ambiente);
			if(ambiente != nullptr){
				response["status"] = "Successful";
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

	
	
	
	
	
	
	
	Json::Value Controller::Crea_Ambiente(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		

		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string cod_ambiente = Current_User ->getemail() + "_" + data["data"]["name"].asString();
		std::string name = data["data"]["name"].asString();

		
		std::string query = "call Ambiente_insert ('"+name +"','"+ Current_User->getemail() + "','"+ cod_ambiente+"');";
		
		if (MYSQL::Query(query) == 0){
				response["data"]["code"] = cod_ambiente;
				response["status"] = "Successful";
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

	Json::Value Controller::Inserisci_Dispositivi(Json::Value data){
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
		std::string cod_ambiente =  data["data"]["envcode"].asString();
		auto entriesArray = data["data"]["devices"];
		Json::Value::iterator devices_to_add;
		std::string start_transaction = "START TRANSACTION;";
		transaction.push_back(start_transaction);
		 for (devices_to_add = entriesArray.begin(); devices_to_add != entriesArray.end();devices_to_add ++){
			std::string name = (*devices_to_add)["name"].asString();
			std::string type = (*devices_to_add)["kind"].asString();
			int code = (*devices_to_add)["code"].asInt();
			
			if((*devices_to_add).isMember("commands")){
				auto commandsEntries = (*devices_to_add)["commands"];
				Json::Value::iterator commands_to_add;
				std::list<std::string> commands;
				std::string query = " insert into Dispositivo (name,type,code,cod_ambiente) values ('"+name+"','"+type+"',"+std::to_string(code)+",'"+ cod_ambiente +"');";
				transaction.push_back(query);
				query = "insert into Attuatore (code) values ("+std::to_string(code) +");";
				transaction.push_back(query);
				for(commands_to_add = commandsEntries.begin(); commands_to_add != commandsEntries.end(); commands_to_add ++ ){
					std::string command = (*commands_to_add).asString();
					query = "insert into Comandi (comando,cod_attuatore) values ('"+ command +"',"+std::to_string(code)+");";
					transaction.push_back(query);
					commands.push_back(command);
				}
				
				Current_User->addDispositivo(cod_ambiente,code,name,type,&commands);
			}
			else{
				
				std::string query = " insert into Dispositivo (name,type,code,cod_ambiente) values ('"+name+"','"+type+"',"+std::to_string(code)+",'"+ cod_ambiente +"');";
				transaction.push_back(query);
				query = "insert into Sensore (code) values ("+std::to_string(code) +");";
				transaction.push_back(query);
				Current_User->addDispositivo(cod_ambiente,code,name,type,nullptr);

			}
			
		}
		std::string commit = "commit;";
		transaction.push_back(commit);
		if (MYSQL::Queries(transaction) ==  false){
			for (devices_to_add = entriesArray.begin(); devices_to_add != entriesArray.end();devices_to_add ++){
				int code = (*devices_to_add)["code"].asInt();
				Current_User->deleteDispositivo(cod_ambiente,code);
			}
			response["status"] = "Failed";
			response["error"] = "Error in creating sensors ";
			response["data"] = "";

		}
		else{
			response["status"] = "Successful";
			response["error"] = "";
			response["data"] = "insert completed";
			std::cout << Current_User->getAmbiente(cod_ambiente)->getDispositivi()->size() << std::endl;

		}
			

	
	
	}

	Json::Value Controller::Visualizza_Ambienti(Json::Value data){
		Json::Value response;

		User * Current_User = Controller::Auth(data["auth"].asString());

		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::list<Ambiente *> * ambienti = Current_User->getAmbienti(); 
		if(ambienti->empty()){
			response["status"] = "Failed";
			response["error"] = "User does not have Ambiente associated";
			response["data"] = "";
			return response;
		}

		std::list<Ambiente *>::iterator ambienti_iterator;
			
		int i = 0;
		for(ambienti_iterator = ambienti->begin();ambienti_iterator != ambienti->end();ambienti_iterator ++) {	
				
			response["data"]["environments"][i] = (*ambienti_iterator)->getcodAmbiente();
			i++;
    	}
        response["status"] = "Successful";
		response["error"]= "";
		return response;
				
	}
			
	Json::Value Controller::Visualizza_Dispositivi(Json::Value data){
		Json::Value response;

		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string cod_ambiente = data["data"]["envcode"].asString();
		int i = 0;
		Ambiente * ambiente = Current_User->getAmbiente(cod_ambiente);
		if(ambiente == nullptr){
			response["status"] = "Failed";
			response["error"] = "Ambiente Not Found";
			response["data"] = "";
			return response;
		}
		std::list<Dispositivo *> * dispositivi = Current_User->getDispositivi(cod_ambiente);
		response["status"] = "Successful";
		response["data"]["devices"]= Json::Value(Json::arrayValue);
		std::list<Dispositivo *>::iterator dispositivi_iterator;
		for(dispositivi_iterator = dispositivi->begin();dispositivi_iterator != dispositivi->end();dispositivi_iterator ++)
		{	
			Json::Value device;
			
			if((*dispositivi_iterator)->get_device_type() == 0 && data["data"]["type"].asString().compare("sensors") == 0){
				device["name"] =(*dispositivi_iterator)->getNome();
				device["kind"] = (*dispositivi_iterator)->getTipo();
				device["code"] = (*dispositivi_iterator)->getCodice();
				response["data"]["devices"][i] = device;
				i++;
			}
			else if((*dispositivi_iterator)->get_device_type() == 1 && data["data"]["type"].asString().compare("actuators")== 0) {
				std::list<std::string>::iterator commands_iterator;
				device["name"] =(*dispositivi_iterator)->getNome();
				device["kind"] = (*dispositivi_iterator)->getTipo();
				device["code"] = (*dispositivi_iterator)->getCodice();
				int x = 0;
				for(commands_iterator = static_cast<Attuatore *>((*dispositivi_iterator))->getComandi()->begin();commands_iterator != static_cast<Attuatore *>((*dispositivi_iterator))->getComandi()->end();commands_iterator ++){
					device["commands"][x] = *commands_iterator;
					x++;
				}
				response["data"]["devices"][i] = device;
				i++;
				
			}
			
			
		}
        response["status"] = "Successful";
				
		return response;
		
	}

	Json::Value Controller::Elimina_Dispositivi(Json::Value data){
		Json::Value response;
		std::list<std::string> transaction;
		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}

		std::string cod_ambiente =  data["data"]["envcode"].asString();
		int i = 0;
		Ambiente * ambiente = Current_User->getAmbiente(cod_ambiente);
		if(ambiente == nullptr){
			response["status"] = "Failed";
			response["error"] = "Ambiente Not Found";
			response["data"] = "";
			return response;
		}
		auto entriesArray = data["data"]["devices"];
		Json::Value::iterator devices_to_delete;
		std::string start_transaction = "START TRANSACTION;";
		transaction.push_back(start_transaction);
		for (devices_to_delete = entriesArray.begin(); devices_to_delete != entriesArray.end();devices_to_delete ++){
			int code_device = (*devices_to_delete).asInt();
			std::string query = "delete from Dispositivo where code = "+std::to_string(code_device) + ";";
			transaction.push_back(query);
		}
		std::string commit = "commit;";
		transaction.push_back(commit);
		if (MYSQL::Queries(transaction) ==  true){
			for (devices_to_delete = entriesArray.begin(); devices_to_delete != entriesArray.end();devices_to_delete ++){
				Controller::User_mutex.lock();
				int cod_sensore = (*devices_to_delete).asInt();
				Current_User->deleteDispositivo(cod_ambiente,cod_sensore);
				Controller::User_mutex.unlock();
			}
			response["status"] = "Successful";
			response["error"] = "";
			response["data"] = "deletion completed";
			return response;

		}
		else{
			response["status"] = "Failed";
			response["error"] = "Error in deleting sensors";
			response["data"] = "";
			return response;
		}

	}

	Json::Value Controller::Invia_Comando(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}

		std::string code_ambiente = data["data"]["envcode"].asString();
		int code = data["data"]["code"].asInt();
		std::string comando = data["data"]["command"].asString();

		if(Current_User->inviaComando(code_ambiente,code,comando)){
			response["status"] = "Successful";
			response["error"] = "";
			response["data"] = "command sent sucessfully";
		}
		else{
			response["status"] = "Failed";
			response["error"] = "Command do not send";
			response["data"] = "";
		}

		return response;

	}

	Json::Value Controller::Visualizza_Storico(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}

		std::string code_ambiente = data["data"]["envcode"].asString();

		std::string query = "select Misure.code,Misure.misura,Misure.time from Misure join Dispositivo on Misure.code = Dispositivo.code where cod_ambiente = '"+code_ambiente+"';";
		
		sql::ResultSet  *res = MYSQL::Select_Query(query);

		if( res->rowsCount() == 0){
				response["status"] = "failed";
				response["error"] = "No data avaible";
				response["data"] = "";
			}
		else{
			response["data"]["sensor_data"]= Json::Value(Json::arrayValue);
			response["status"] = "Successful";
			response["error"] = "";
			int i = 0;
			while (res->next()) {
				Json::Value row;
				row["code"] = res->getInt("code");
				row["misura"] = (float)res->getDouble("misura");
				row["time"] = (std::string)res->getString("time");
				response["data"]["sensor_data"][i]=row;
				i++;
			}
			return response;

		}
	}
	Json::Value Controller::Elimina_Ambiente(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string code_ambiente =  data["data"]["envcode"].asString();
		std::string query = "delete from Ambiente where cod_ambiente = '" + code_ambiente +"'";
		if(Current_User->eliminaAmbiente(code_ambiente) && MYSQL::Query(query) == 0){
			MQTTClient_unsubscribe(client, code_ambiente.c_str());
			response["status"] = "Successful";
			response["error"] = "";
			response["data"] = "Ambiente "+ data["data"]["encode"].asString() +" deleted";
		}
		else {
			response["status"] = "Failed";
			response["error"] = "Ambiente "+ data["data"]["envcode"].asString() +" not found";
			response["data"] = "";
		}

		return response;


	}

	Json::Value Controller::Logout(Json::Value data){
		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr){
			response["status"] = "Failed";
			response["error"] = "No User Found";
			response["data"] = "";
			return response;
		}


		users.remove(Current_User);
		delete Current_User;
		

		response["status"] = "Successful";
		response["error"] = "";
		response["data"] = "Logout";
		return response;

	}

	Controller::~Controller(){
		User_mutex.lock();
		std::list<User *>::iterator users_iterator;
		for(User * user : users){
			delete user;
		}
        users.clear();
		User_mutex.unlock();
		instance = 0;
		MQTTClient_disconnect(client, 10000);
        MQTTClient_destroy(&client);

	}

	Json::Value Controller::Associa_Utente(Json::Value data){

		Json::Value response;
		User * Current_User = Controller::Auth(data["auth"].asString());
		
		if(Current_User == nullptr || Current_User->getAdmin() == false){
			response["status"] = "Failed";
			response["error"] = "Unauthorized";
			response["data"] = "";
			return response;
		}
		std::string code_ambiente =  data["data"]["envcode"].asString();
		Ambiente * ambiente = Current_User->getAmbiente(code_ambiente);
		if(ambiente == nullptr){
			response["status"] = "Failed";
			response["error"] = "Ambiente not found";
			response["data"] = "";
			return response;
		}
		std::string query = "insert into Ambiente_User (cod_ambiente,User_email) values ('"+ code_ambiente +"','"+data["data"]["user"].asString()+"');";
		if(MYSQL::Query(query) != 0){
			response["status"] = "Failed";
			response["error"] = "User not found or ambiente already assiociated";
			response["data"] = "";
			return response;
		}

		for(User * user : users){
			if (user->getemail().compare(data["data"]["user"].asString()) == 0){
				if(user->getAmbiente(code_ambiente) == nullptr){
					user->getAmbienti()->push_back(ambiente);
				}
			}
		}
		response["status"] = "Successful";
		response["error"] = "";
		response["data"] = "Ambiente associated";
		return response;
		
		}

		
		

	

