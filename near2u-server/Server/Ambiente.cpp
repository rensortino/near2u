

#include "Ambiente.hpp"


    Ambiente::Ambiente(std::string& name, std::string& cod_ambiente){
        Ambiente::Nome = name;
        Ambiente::codAmbiente = cod_ambiente;
    }
    std::string& Ambiente::getcodAmbiente(){
        return Ambiente::codAmbiente;
    }
    std::string& Ambiente::getNome(){
        return Ambiente::Nome;
    }
    std::list<Dispositivo *> * Ambiente::getDispositivi(){
        return &dispositivi;
    }


    void Ambiente::deleteDispositivo(int cod_dispositivo){

         Dispositivo * dispositivo = getDispositivo(cod_dispositivo);
         delete dispositivo;
         dispositivi.remove(dispositivo);
            
        
    }

    void Ambiente::addDispositivo(int code, std::string& nome, std::string& tipo, std::list<std::string> * commands){

        if(commands != nullptr){
            std::cout << "inserisco attuatore" << std::endl;
            Attuatore *  dispositivo = new Attuatore(code,nome,tipo,commands);
            dispositivi.push_back(dispositivo);
        }
        else {
            std::cout << "inserisco sensore" << std::endl;
            Sensore *  dispositivo = new Sensore(code,nome,tipo);
            dispositivi.push_back(dispositivo);
        }
       
        
    }

    Dispositivo * Ambiente::getDispositivo(int code){

        for(Dispositivo * dispositivo : dispositivi){
            if(dispositivo->getCodice() == code){
                return dispositivo;
                
            }
        }

        return nullptr;


    }

    bool Ambiente::inviaComando(int code_attuatore, std::string& comando){
        Attuatore * attuatore;
        try{
            attuatore = static_cast<Attuatore *>(getDispositivo(code_attuatore));
        }
        catch(std::exception &error){
            std::cout << error.what() << std::endl;
            return false;
        }

        if(attuatore->controllaComando(comando) == false){
            return false;
        }

        std::string topic = codAmbiente + std::to_string(code_attuatore);
        std::cout << topic << std::endl;

        MQTTClient client;
        MQTTClient_connectOptions conn_opts = MQTTClient_connectOptions_initializer;
        MQTTClient_message pubmsg = MQTTClient_message_initializer;
        MQTTClient_deliveryToken token;
        int rc;

        MQTTClient_create(&client, ADDRESS, CLIENTSERVER,
            MQTTCLIENT_PERSISTENCE_NONE, NULL);
            conn_opts.keepAliveInterval = 20;
            conn_opts.cleansession = 1;

        if ((rc = MQTTClient_connect(client, &conn_opts)) != MQTTCLIENT_SUCCESS)
        {
            printf("Failed to connect, return code %d\n", rc);
            return false;
        }

        std::string message = "{\"code\":" + std::to_string(attuatore->getCodice()) + ",\"command\":\""+ comando + "\"}";
        std::cout << message << std::endl;
        pubmsg.payload = (void *)message.c_str();
        pubmsg.payloadlen = message.size();
        pubmsg.qos = QOS;
        pubmsg.retained = 0;
        MQTTClient_publishMessage(client, topic.c_str(), &pubmsg, &token);
        rc = MQTTClient_waitForCompletion(client, token, TIMEOUT);
        printf("Command with delivery token %d delivered\n", token);

        
        MQTTClient_disconnect(client, 10000);
        MQTTClient_destroy(&client);

        return true;
        


    }

    Ambiente::~Ambiente(){
     
        for(Dispositivo * dispositivo : dispositivi){
            delete dispositivo;
        }
        dispositivi.clear();
    }
    Sensore * Ambiente::getSensore(int cod_sensore){
        std::list<Sensore>::iterator sensori_iterator;

        for(sensori_iterator=sensori.begin(); sensori_iterator != sensori.end(); sensori_iterator ++){
            if((*sensori_iterator).getCodSensore() == cod_sensore){
                return &(*sensori_iterator);
            }
        }
        return nullptr;


    }


  
