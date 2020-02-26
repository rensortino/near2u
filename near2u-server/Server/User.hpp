#ifndef USER_HPP
#define USER_HPP
#include <string>
#include <iostream>
#include <assert.h>
#include "Ambiente.hpp" 


class User
{

public:
     User(const std::string& n, const std::string& s ,const std::string& e , const std::string& p,const std::string& pa);
     std::string getName();
     std::string getsurname();
     std::string getemail();
     std::string getPassword();
     std::string getauth_token();
     std::list<Ambiente> * getAmbienti();
     bool getAdmin();
     void setAdmin(bool role);
    Ambiente * getAmbiente(std::string cod_Ambiente);
private:
    std::string name;
    std::string surname;
    std::string email;
    std::string password;
    std::string auth_token;
    bool admin;
 
    std::list<Ambiente> ambienti;
	

};

#endif
 
