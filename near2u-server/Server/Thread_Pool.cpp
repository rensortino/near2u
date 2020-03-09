
#include "Thread_Pool.hpp"

<<<<<<< HEAD
    enum StringValue { Default,
                        Register, 
                        Login, 
                        Topic_Ambiente,
                        Configura_Ambiente,
                        Inserisci_Sensori,
                        Visualizza_Ambienti,
                        Visualizza_Sensori,
                        Elimina_Sensori
                        };
    static std::map<std::string, StringValue> s_mapStringValues;
=======
    
   
>>>>>>> Iterazione_3

    static void Initialize();

    Thread_Pool::Thread_Pool() {
        
        done = false;
        Initialize();
        auto numberOfThreads = std::thread::hardware_concurrency();
        if (numberOfThreads == 0) {
        numberOfThreads = 1;
        }
        unsigned i ;
        for(i = 0 ; i < numberOfThreads; i ++){
            threads.push_back(std::thread(&Thread_Pool::TaskWork,this));
        }
    }

    void Thread_Pool::stop(){
        done = true; 
        workQueueConditionVariable.notify_all();
        for(auto& thread : threads){
            if(thread.joinable()){
                thread.join();
            }
        }
        delete controller;
        std::cout << "Thread pool deleted" <<std::endl;
    }

    void Thread_Pool::queueWork(int fd, std::string& request){

        std::lock_guard<std::mutex> g(QueueMutex);
        requestqueue.push(std::pair<int, std::string>(fd, request));
        workQueueConditionVariable.notify_one();
    }

    void Thread_Pool::TaskWork(){
        while(!done){
            std::pair<int,std::string> request;
            {
            std::unique_lock<std::mutex> g(QueueMutex);
            workQueueConditionVariable.wait(g, [&]{
            return !requestqueue.empty() || done;
            });
                request = requestqueue.front();
                requestqueue.pop();
                ElaborateRequest(request);  
            }
            
        }
    }

    void Thread_Pool::ElaborateRequest(const std::pair<int,std::string> request){
        Json::Reader reader;
        Json::Value requestjson;
        std::string response;
        controller = Controller::getIstance();
        controller->setUpMqtt();
        reader.parse(request.second, requestjson);
        std::cout << "new request arrived requesting API: " + requestjson["function"].asString() <<std::endl;
        switch (s_mapStringValues[requestjson["function"].asString()])
        {
            case Register:
                response =  controller->Register(requestjson["data"]).toStyledString(); 
                break;
            case Login:
                response =  controller->Login(requestjson["data"]).toStyledString();  
                break;
            case Topic_Ambiente:
                response =  controller->Topic_Ambiente(requestjson).toStyledString(); 
                break;
            
            case Crea_Ambiente:
                response = controller->Crea_Ambiente(requestjson).toStyledString();
                break;
            case Inseresci_Dispositivi:
                response = controller->Inserisci_Dispositivi(requestjson).toStyledString();
                break;
            case Visualizza_Ambienti:
                response = controller->Visualizza_Ambienti(requestjson).toStyledString();
                break;
            case Visualizza_Dispositivi:
                response = controller->Visualizza_Dispositivi(requestjson).toStyledString();
                break;
            case Elimina_Dispositivi:
                response = controller->Elimina_Dispositivi(requestjson).toStyledString();
                break;
            case Invia_Comando:
                response = controller->Invia_Comando(requestjson).toStyledString();
                break;
            case Visualizza_Storico:
                response = controller->Visualizza_Storico(requestjson).toStyledString();
                break;
            case Elimina_Ambiente:
                response = controller->Elimina_Ambiente(requestjson).toStyledString();
                break;
            case Logout:
                response = controller->Logout(requestjson).toStyledString();
                break;
            case Associa_Utente:
                response = controller->Associa_Utente(requestjson).toStyledString();
                break;
            case Visualizza_Ambienti:
                response = controller->Visualizza_Ambienti(requestjson).toStyledString();
                break;
            case Visualizza_Sensori:
                response = controller->Visualizza_Sensori(requestjson).toStyledString();
                break;
                case Elimina_Sensori:
                response = controller->Elimina_sensori(requestjson).toStyledString();
                break;
            default:
                response = "{\"status\" : \"Service not avaible\"}"; 
                break;
        }
        send(request.first, response.c_str(), response.size(), 0);
        // Close the connection
        close(request.first);
    }



    void Thread_Pool::Initialize(){
    s_mapStringValues["register"] = Register;
    s_mapStringValues["login"] = Login;
    s_mapStringValues["topic_ambiente"] = Topic_Ambiente;
<<<<<<< HEAD
    s_mapStringValues["configura_ambiente"] = Configura_Ambiente;
    s_mapStringValues["inserisci_sensori"] = Inserisci_Sensori;
    s_mapStringValues["visualizza_ambienti"] = Visualizza_Ambienti;
    s_mapStringValues["visualizza_sensori"] = Visualizza_Sensori;
    s_mapStringValues["elimina_sensori"] = Elimina_Sensori;
    
    std::cout << "s_mapStringValues contains " 
        << s_mapStringValues.size() 
        << " entries." << std::endl;
=======
    s_mapStringValues["crea_ambiente"] = Crea_Ambiente;
    s_mapStringValues["inserisci_dispositivi"] = Inseresci_Dispositivi;
    s_mapStringValues["visualizza_ambienti"] = Visualizza_Ambienti;
    s_mapStringValues["elimina_dispositivi"] = Elimina_Dispositivi;
    s_mapStringValues["visualizza_dispositivi"] = Visualizza_Dispositivi;
    s_mapStringValues["invia_comando"] = Invia_Comando;
    s_mapStringValues["visualizza_storico"] = Visualizza_Storico;
    s_mapStringValues["elimina_ambiente"] = Elimina_Ambiente;
    s_mapStringValues["logout"] = Logout;
    s_mapStringValues["associa_utente"] = Associa_Utente;
>>>>>>> Iterazione_3
    }
