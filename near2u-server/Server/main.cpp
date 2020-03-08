#include <sys/socket.h> // For socket()
#include <netinet/in.h> // For sockaddr_in
#include <thread>
#include "Thread_Pool.hpp"
#include "Ambiente_Simulation.cpp"


void stop_server( bool& done){
    char ch;
    printf("Press Q<Enter> to quit\n\n");
    do 
      {
        ch = getchar();
      } 
    while(ch!='Q' && ch != 'q');

    done = true;
    
    }


int main(){
    auto tp =   new Thread_Pool();
    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
    // creo un thread che simula la pubblicazione dei vari sensori per ogni ambiente su MQTT
    bool done = false;
    std::thread simulation(sensors_pubblish,std::ref(done));
    std::thread stop(stop_server,std::ref(done));
    
    
  if (sockfd == 0) {
    std::cout << "Failed to create socket. errno: " << errno << std::endl;
    exit(EXIT_FAILURE);
  }

   sockaddr_in sockaddr;
  sockaddr.sin_family = AF_INET;
  sockaddr.sin_addr.s_addr = INADDR_ANY;
  sockaddr.sin_port = htons(3333);

  if (bind(sockfd, (struct sockaddr*)&sockaddr, sizeof(sockaddr)) < 0) {
    std::cout << "Failed to bind to port 333. errno: " << errno << std::endl;
    exit(EXIT_FAILURE);
  }

  if (listen(sockfd, 10) < 0) {
    std::cout << "Failed to listen on socket. errno: " << errno << std::endl;
    exit(EXIT_FAILURE);
  }

  while(!done){
    auto addrlen = sizeof(sockaddr);
    int connection = accept(sockfd, (struct sockaddr*)&sockaddr, (socklen_t*)&addrlen);
    if (connection < 0) {
      std::cout << "Failed to grab connection. errno: " << errno << std::endl;
      exit(EXIT_FAILURE);
    }

     // Read from the connection
    char buffer[8192];
    auto bytesRead = read(connection, buffer, 8192);
    std::string request = buffer;

    // Add some work to the queue
    tp->queueWork(connection, request);
  }

  
  
  simulation.join();
  stop.join();
  tp->stop();
  close(sockfd);
  printf("server shutting down\n");
  



}