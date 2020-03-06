#include "function_mqtt.hpp"

namespace MQTT{
    void delivered(void *context,MQTTClient_deliveryToken token){
        printf("Message delivery confirmed\n");
    }

    int UploadDataSensor(void *context, char *topicName, int topiclen, MQTTClient_message *message){
        Json::Value sensor_data;
        Json::Reader reader;
        int i;
        char* payloadptr;


        payloadptr =(char*) message->payload;
        reader.parse(payloadptr,sensor_data);

        std::string query = "insert into Misure (misura,code,time) values ("+ std::to_string(sensor_data["measurement"].asFloat())+","+std::to_string(sensor_data["code"].asInt())+",'"+sensor_data["time"].asString() +"');";
        std::cout << query << std::endl;
        MYSQL::Query(query);


        MQTTClient_freeMessage(&message);
        MQTTClient_free(topicName);
        return 1;
    }

    void connlost(void *context, char *cause){
        printf("\nConnection lost\n");
        printf("     cause: %s\n", cause);
    }

    int msgarrvd(void *context, char *topicName, int topicLen, MQTTClient_message *message)
    {
        int i;
        char* payloadptr;

        printf("Message arrived\n");
        printf("     topic: %s\n", topicName);
        printf("   message: ");
        printf("%s\n",(char*) message->payload);

        
        MQTTClient_freeMessage(&message);
        MQTTClient_free(topicName);
        return 1;
    }

    MQTTClient connect_publisher(std::string& address, std::string& ClientId){

        MQTTClient client;
        MQTTClient_connectOptions conn_opts = MQTTClient_connectOptions_initializer;
        int rc;

        MQTTClient_create(&client, address.c_str(), ClientId.c_str(),
            MQTTCLIENT_PERSISTENCE_NONE, NULL);
        conn_opts.keepAliveInterval = 20;
        conn_opts.cleansession = 1;

        if ((rc = MQTTClient_connect(client, &conn_opts)) != MQTTCLIENT_SUCCESS)
        {
            printf("Failed to connect, return code %d\n", rc);
            return nullptr;
        }
        return client;

    }

    MQTTClient connect_subscriber(std::string& address, std::string& ClientId,int func_on_receive){

        MQTTClient client;
        MQTTClient_connectOptions conn_opts = MQTTClient_connectOptions_initializer;
        int rc;

        MQTTClient_create(&client, address.c_str(), ClientId.c_str(),
            MQTTCLIENT_PERSISTENCE_NONE, NULL);
        conn_opts.keepAliveInterval = 20;
        conn_opts.cleansession = 1;

        if ((rc = MQTTClient_connect(client, &conn_opts)) != MQTTCLIENT_SUCCESS)
        {
            printf("Failed to connect, return code %d\n", rc);
            return nullptr;
        }
        if (func_on_receive == 0){
            MQTTClient_setCallbacks(client, NULL, connlost, msgarrvd, delivered);
        }
        else {
            MQTTClient_setCallbacks(client, NULL, connlost, UploadDataSensor, delivered);
        }

        
        return client;

    }
    bool publish(std::string& topic, std::string& message, MQTTClient client ){
        MQTTClient_message pubmsg = MQTTClient_message_initializer;
        MQTTClient_deliveryToken token;
        int rc;
        pubmsg.payload = (void *)(message.c_str());
        pubmsg.payloadlen = message.size();
        pubmsg.qos = QOS;
        pubmsg.retained = 0;
        MQTTClient_publishMessage(client, topic.c_str(), &pubmsg, &token);
        rc = MQTTClient_waitForCompletion(client, token, TIMEOUT);
        if(rc == MQTTCLIENT_SUCCESS){
            printf("Message with delivery token %d delivered\n", token);
            return true;
        }
        else {
             printf("Error in delivering message");
             return false;
        }
    }
    void subscribe(std::string& topic, MQTTClient client){
        MQTTClient_subscribe(client, topic.c_str(), QOS);
        
    }
}