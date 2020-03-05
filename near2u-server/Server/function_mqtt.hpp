#include "MQTTClient.h"
#include <jsoncpp/json/json.h>
#include "MYSQL.hpp"
#include <string>
#include <iostream>

void delivered(void *context,MQTTClient_deliveryToken token);

int UploadDataSensor(void *context, char *topicName, int topiclen, MQTTClient_message *message);

void connlost(void *context, char *cause);
int msgarrvd(void *context, char *topicName, int topicLen, MQTTClient_message *message);
