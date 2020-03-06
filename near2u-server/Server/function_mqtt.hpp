#include "MQTTClient.h"
#include <jsoncpp/json/json.h>
#include "MYSQL.hpp"
#include <string>
#include <iostream>

namespace MQTT{

void delivered(void *context,MQTTClient_deliveryToken token);

int UploadDataSensor(void *context, char *topicName, int topiclen, MQTTClient_message *message);

void connlost(void *context, char *cause);
int msgarrvd(void *context, char *topicName, int topicLen, MQTTClient_message *message);

MQTTClient connect(std::string& address, std::string& ClientId);
MQTTClient connect_subscriber(std::string& address, std::string& ClientId);
}
