#include "Controller.hpp"
#include "Ambiente.hpp"
#include "Sensore.h"
#include <iostream> 
#include <thread>         
#include <chrono>
#include "stdio.h"
#include "stdlib.h"
#include "string.h"
#include "function_mqtt.hpp"
#include <ctime>
#include <time.h>



// questa funzione serve a simulare la pubblicazione sul broker mqtt da parte dei dispositivi dei vari ambienti
using namespace MQTT;

void sensors_pubblish(){
    std::string address(getenv("MQTT_BROKER_ADDRESS"));
    std::string sensor_id("sensors_simulation");
    std::string actuator_id("actuators_simulation");

    MQTTClient client = connect(address,sensor_id);

    MQTTClient client_attuatori = connect_subscriber(address,actuator_id,0);
    std::list<std::string> lista_topic_attuatori;
    std::list<std::string>::iterator topic_iterator;

    Controller * controller = Controller::getIstance();

    while(true){
        std::this_thread::sleep_for (std::chrono::seconds(5));
        std::cout<< "checking for sensors" << std::endl;
        std::list<User>::iterator users_iterator;

        controller->getUser_mutex()->lock_shared();
        for(users_iterator = controller->getUsers()->begin(); users_iterator != controller->getUsers()->end(); users_iterator ++ ){
            std::list<Ambiente>::iterator ambienti_itarator;
            std::list<Ambiente> * ambienti = (*users_iterator).getAmbienti(); 
            for(ambienti_itarator = ambienti->begin(); ambienti_itarator != ambienti->end(); ambienti_itarator ++){
                std::string topic = (*ambienti_itarator).getcodAmbiente();
                std::list<Dispositivo *>::iterator dispositivi_iterator;
                std::list<Dispositivo *> * dispositivi = (*ambienti_itarator).getDispositivi(); 
                for(dispositivi_iterator = dispositivi->begin();dispositivi_iterator != dispositivi->end(); dispositivi_iterator ++){
                    if((*dispositivi_iterator)->get_device_type() == device_type::sensore){
                        std::time_t result = std::time(nullptr);
                        time_t t = time(NULL);
                        struct tm tm = *localtime(&t);

                        printf("%d/%d/%d %d:%d",tm.tm_year,tm.tm_mon,tm.tm_mday,tm.tm_hour,tm.tm_sec);
                        
                        std::string message = "{\"code\":" + std::to_string((*dispositivi_iterator)->getCodice()) + ",\"name\":\""+ (*dispositivi_iterator)->getNome() + "\",\"type\":\""+(*dispositivi_iterator)->getTipo() +" \",\"measurement\":"+std::to_string((float)rand()/(float)(RAND_MAX/15)) +",\"time\": \""+ std::asctime(std::localtime(&result)) +"\" }";
                        publish(topic,message,client);
                    }
                    else if((*dispositivi_iterator)->get_device_type() == device_type::attuatore){
                       std::string topic_attuatore = (*ambienti_itarator).getcodAmbiente() + std::to_string((*dispositivi_iterator)->getCodice()); 
                       if(lista_topic_attuatori.empty()){
                           std::cout << "attuatore: " + std::to_string((*dispositivi_iterator)->getCodice()) + "subscribing to: " + topic_attuatore   << std::endl;
                            subscribe(topic,client_attuatori);
                       }
                       for(topic_iterator = lista_topic_attuatori.begin(); topic_iterator != lista_topic_attuatori.end(); topic_iterator ++ ){
                            if((*topic_iterator).compare(topic_attuatore) == 0){
                                break;
                            }
                            else{
                                std::cout << "attuatore: " + std::to_string((*dispositivi_iterator)->getCodice()) + "subscribing to: " + topic_attuatore   << std::endl;
                                subscribe(topic,client_attuatori);

                            }
                        }
                        lista_topic_attuatori.push_back(topic_attuatore);
                       
                    }
                }
            }
        }
        controller->getUser_mutex()->unlock_shared();


        
    }
    
    for(topic_iterator = lista_topic_attuatori.begin(); topic_iterator != lista_topic_attuatori.end(); topic_iterator ++ ){
        MQTTClient_unsubscribe(client, (*topic_iterator).c_str());
    }
    MQTTClient_disconnect(client_attuatori, 10000);
    MQTTClient_destroy(&client_attuatori);
    MQTTClient_disconnect(client, 10000);
    MQTTClient_destroy(&client);
    exit(0);
}