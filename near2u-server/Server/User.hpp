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
     ~User();
     std::string& getName();
     std::string& getsurname();
     std::string& getemail();
     std::string& getPassword();
     std::string& getauth_token();
     std::list<Ambiente *> * getAmbienti();
     bool getAdmin();
     void setAdmin(bool role);
    Ambiente * getAmbiente(std::string& cod_Ambiente);
    void addAmbiente(std::string& nome, std::string& codice );
    std::list<Dispositivo *> * getDispositivi(std::string& code);
    void deleteDispositivo(std::string& cod_ambiente,int code);
    void addDispositivo(std::string& cod_ambiente,int code, std::string& nome, std::string& tipo, std::list<std::string>* commands);
    bool inviaComando(std::string& cod_ambiente, int code_attuatore, std::string& comando);
    bool eliminaAmbiente(std::string& cod_ambiente);
private:
    std::string name;
    std::string surname;
    std::string email;
    std::string password;
    std::string auth_token;
    bool admin;
 
    std::list<Ambiente *> ambienti;
	

};

#endif
 
