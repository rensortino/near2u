
#include "Thread_Pool.hpp"
#include "Controller.hpp"

    enum StringValue { Default,
                        Register, 
                        Login, 
                        Topic_Ambiente,
                        Configura_Ambiente,
                        Inseresci_Dispositivi,
                        Visualizza_Ambienti,
                        Visualizza_Sensori,
                        Visualizza_Dispositivi,
                        Elimina_Sensori
                        };
    static std::map<std::string, StringValue> s_mapStringValues;

    static void Initialize();

    Thread_Pool::Thread_Pool() {
        done = false;
        Initialize();
        // set the number of thread depending on the hardware if the hardware is not multithread set the number of thread to 1
        auto numberOfThreads = std::thread::hardware_concurrency();
        if (numberOfThreads == 0) {
        numberOfThreads = 1;
        }
        for(unsigned i = 0 ; i < numberOfThreads; i ++){
            threads.push_back(std::thread(&Thread_Pool::TaskWork,this));
            // we populate the thread vector indicating the function each thread should execute
        }
    }
    Thread_Pool::~Thread_Pool(){
        done = true; //Indicate that the server is shutting down
        workQueueConditionVariable.notify_all();
        for(auto& thread : threads){
            if(thread.joinable()){
                thread.join();
            }
        }
    }

    void Thread_Pool::queueWork(int fd /* file descriptor for socket */, std::string& request){

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
            // Only wake up if there are elements in the queue or the program is
            // shutting down
            return !requestqueue.empty() || done;
            });

            request = requestqueue.front();
            requestqueue.pop();
        }
        ElaborateRequest(request);
        }
    }

    void Thread_Pool::ElaborateRequest(const std::pair<int,std::string> request){
        Json::Reader reader;
        Json::Value requestjson;
        std::string response;
        Controller * controller = Controller::getIstance();
        reader.parse(request.second, requestjson);
        std::cout << "new request arrived requesting API: " + requestjson["function"].asString() <<std::endl;
        std::cout << s_mapStringValues[requestjson["function"].asString()] << std::endl;
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
            
            case Configura_Ambiente:
                response = controller->Configura_ambiente(requestjson).toStyledString();
                break;
            case Inseresci_Dispositivi:
                response = controller->Inserisci_Dispositivi(requestjson).toStyledString();
                break;
            case Visualizza_Ambienti:
                response = controller->Visualizza_Ambienti(requestjson).toStyledString();
                break;
            case Visualizza_Sensori:
               // response = controller->Visualizza_Sensori(requestjson).toStyledString();
                break;
                case Elimina_Sensori:
                //response = controller->Elimina_sensori(requestjson).toStyledString();
                break;
                case Visualizza_Dispositivi:
                response = controller->Visualizza_Dispositivi(requestjson).toStyledString();
                break;
            default:
                response = "{\"status\" : \"Service not avaible\"}"; 
                break;
        }
        send(request.first, response.c_str(), response.size(), 0);
        // Close the connection
        close(request.first);
    }


    void Initialize()
    {
    s_mapStringValues["register"] = Register;
    s_mapStringValues["login"] = Login;
    s_mapStringValues["topic_ambiente"] = Topic_Ambiente;
    s_mapStringValues["configura_ambiente"] = Configura_Ambiente;
    s_mapStringValues["inserisci_dispositivi"] = Inseresci_Dispositivi;
    s_mapStringValues["visualizza_ambienti"] = Visualizza_Ambienti;
    s_mapStringValues["visualizza_sensori"] = Visualizza_Sensori;
    s_mapStringValues["elimina_sensori"] = Elimina_Sensori;
    s_mapStringValues["visualizza_dispositivi"] = Visualizza_Dispositivi;
    
    std::cout << "s_mapStringValues contains " 
        << s_mapStringValues.size() 
        << " entries." << std::endl;
    }
