
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


    class Thread_Pool {

    public:
        Thread_Pool();
        ~Thread_Pool();
        void queueWork(int fd /* file descriptor for socket */, std::string& request);
        

    private:
        std::condition_variable_any workQueueConditionVariable;

        std::vector<std::thread> threads; // We store the threads in this variable

        //we need a mutex for accessing che queue
        std::mutex QueueMutex;

        //we need a queue of request that has to be processed

        std::queue<std::pair<int,std::string>> requestqueue;

        bool done; // to notify to stop the server

        void TaskWork();
        void ElaborateRequest(const std::pair<int,std::string>);

    };