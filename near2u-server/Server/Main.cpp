#include "Controller.cpp"
#include <unistd.h> 
#include <sys/socket.h> 
#include <stdlib.h> 
#include <netinet/in.h> 
#include <stdio.h> 
#include <string.h>
#include <fcntl.h> 
#include<pthread.h>


#define HOST "localhost"
#define USER "admin"
#define PASS "admin"
#define DB "apl_project"
#define PORT 3333 



using namespace Server;



int main(){
    Json::Reader reader;
    Json::Value request;
    int server_fd, new_socket, valread; 
    struct sockaddr_in address; 
    int opt = 1; 
    int addrlen = sizeof(address); 
    char buffer[8192] = {0}; 
    std::string response;
    Controller * controller = Controller::getIstance();
     

    std::cout << "Server is Starting" <<std::endl;
    // Creating socket file descriptor 
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) 
    { 
        perror("socket failed"); 
        exit(EXIT_FAILURE); 
    } 
       
    // Forcefully attaching socket to the port 8080 
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, 
                                                  &opt, sizeof(opt))) 
    { 
        perror("setsockopt"); 
        exit(EXIT_FAILURE); 
    } 
    address.sin_family = AF_INET; 
    address.sin_addr.s_addr = INADDR_ANY; 
    address.sin_port = htons( PORT ); 
       
    // Forcefully attaching socket to the port 8080 
    if (bind(server_fd, (struct sockaddr *)&address,  
                                 sizeof(address))<0) 
    { 
        perror("bind failed"); 
        exit(EXIT_FAILURE); 
    } 
    
        if (listen(server_fd, 3) < 0) 
        { 
            perror("listen"); 
            exit(EXIT_FAILURE); 
        } 
        if ((new_socket = accept(server_fd, (struct sockaddr *)&address,  
                        (socklen_t*)&addrlen))<0) 
        { 
            perror("accept"); 
            exit(EXIT_FAILURE); 
        } 
    while(1){   
        valread = read( new_socket , buffer, 8192); 
        reader.parse(buffer, request);
        std::cout << "new request arrived requesting API: " + request["function"].asString() <<std::endl;
        if(request["function"].asString().compare("register") == 0){
          response =  controller->Register(request["data"]).toStyledString(); 
        }
        if(request["function"].asString().compare("login") == 0){
          response =  controller->Login(request["data"]).toStyledString(); 
        }
        if(request["function"].asString().compare("seleziona ambiente") == 0){
          response =  controller->Seleziona_Ambiente(request["data"]).toStyledString(); 
        }
        char response_ctr[response.size()+ 1];
        strcpy(response_ctr,response.c_str());
        send(new_socket,response_ctr,strlen(response_ctr),0);
        
        
    }
    return 0 ;
}