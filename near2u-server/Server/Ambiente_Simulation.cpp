#include "Controller.hpp"
#include "Ambiente.hpp"
#include "Sensore.h"
#include <iostream> 
#include <thread>         
#include <chrono>
#include "stdio.h"
#include "stdlib.h"
#include "string.h"
#include "MQTTClient.h"

#define ADDRESS     "tcp://localhost:8082"
#define CLIENTID    "ambienti_simulation"
#define QOS         1
#define TIMEOUT     10000L


// questa funzione serve a simulare la pubblicazione sul broker mqtt da parte dei sensori dei vari ambienti
void sensors_pubblish(){

    MQTTClient client;
    MQTTClient_connectOptions conn_opts = MQTTClient_connectOptions_initializer;
    MQTTClient_message pubmsg = MQTTClient_message_initializer;
    MQTTClient_deliveryToken token;
    int rc;

    MQTTClient_create(&client, ADDRESS, CLIENTID,
        MQTTCLIENT_PERSISTENCE_NONE, NULL);
    conn_opts.keepAliveInterval = 20;
    conn_opts.cleansession = 1;

    if ((rc = MQTTClient_connect(client, &conn_opts)) != MQTTCLIENT_SUCCESS)
    {
        printf("Failed to connect, return code %d\n", rc);
        exit(-1);
    }

    Controller * controller = Controller::getIstance();

    while(true){
        std::this_thread::sleep_for (std::chrono::seconds(5));

        std::list<User>::iterator users_iterator;

        controller->getUser_mutex()->lock_shared();
        for(users_iterator = controller->getUsers()->begin(); users_iterator != controller->getUsers()->end(); users_iterator ++ ){
            std::list<Ambiente>::iterator ambienti_itarator;
            std::list<Ambiente> * ambienti = (*users_iterator).getAmbienti(); 
            for(ambienti_itarator = ambienti->begin(); ambienti_itarator != ambienti->end(); ambienti_itarator ++){
                std::string topic = (*ambienti_itarator).getcodAmbiente();
                std::list<Sensore>::iterator sensori_itarator;
                std::list<Sensore> * sensori = (*ambienti_itarator).getSensori(); 
                for(sensori_itarator = sensori->begin();sensori_itarator != sensori->end(); sensori_itarator ++){
                    std::string message = "{\"code\":" + std::to_string(sensori_itarator->getCodSensore()) + ",\"name\":\""+ sensori_itarator->getName() + "\",\"type\":\""+sensori_itarator->getType() +" \",\"measurement\":"+std::to_string((float)rand()/(float)(RAND_MAX/15)) +" }";
                    std::cout << message << std::endl;
                    pubmsg.payload = &message;
                    pubmsg.payloadlen = message.size();
                    pubmsg.qos = QOS;
                    pubmsg.retained = 0;
                    MQTTClient_publishMessage(client, topic.c_str(), &pubmsg, &token);
                    rc = MQTTClient_waitForCompletion(client, token, TIMEOUT);
                    printf("Message with delivery token %d delivered\n", token);
                }
            }
        }
        controller->getUser_mutex()->unlock_shared();


        
    }
    MQTTClient_disconnect(client, 10000);
    MQTTClient_destroy(&client);
    exit(0);
}