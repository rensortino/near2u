#ifndef USER_H
#define USER_H
#include <string>
#include <iostream>
#include <assert.h>
#include "Ambiente.h" 

namespace Server
{
class User
{

public:
     User(const std::string& n, const std::string& s ,const std::string& e , const std::string& p);
     std::string getName();
     std::string getsurname();
     std::string getemail();
     std::string getauth_token();
     std::list<Ambiente> * getAmbienti();
private:
    std::string name;
    std::string surname;
    std::string email;
    std::string password;
    std::string auth_token;
 
    std::list<Ambiente> ambienti;
	

};

} 
#endif
 
