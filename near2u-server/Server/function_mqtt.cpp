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

}