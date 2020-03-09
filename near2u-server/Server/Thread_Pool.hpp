
#include <sys/socket.h>
#include <cstdlib> // For exit() and EXIT_FAILURE
#include <iostream> // For cout
#include <unistd.h> // For read
#include <thread> // std::thread
#include <vector> // std::vector
#include <queue> // std::queue
#include <mutex> // std::mutex
#include <map>
#include <condition_variable> // std::condition_variable
#include "Controller.hpp"

enum StringValue {      Default,
                        Register, 
                        Login, 
                        Topic_Ambiente,
                        Crea_Ambiente,
                        Inseresci_Dispositivi,
                        Visualizza_Ambienti,
                        Visualizza_Dispositivi,
                        Elimina_Dispositivi,
                        Invia_Comando,
                        Visualizza_Storico,
                        Elimina_Ambiente,
                        Logout,
                        Associa_Utente,
                        };


    class Thread_Pool {

    public:
        Thread_Pool();
        void queueWork(int fd /* file descriptor for socket */, std::string& request);
        void Initialize();
        void stop();
        
        Controller * controller;
    private:
        std::condition_variable_any workQueueConditionVariable;
        std::map<std::string, StringValue> s_mapStringValues;
        std::vector<std::thread> threads;
        std::mutex QueueMutex;
        std::queue<std::pair<int,std::string>> requestqueue;
        bool done; 
        void TaskWork();
        void ElaborateRequest(const std::pair<int,std::string>);
    };