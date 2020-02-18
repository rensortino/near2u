
#include "User.h"

namespace Server
{
    User::User(const std::string& n, const std::string& s ,const std::string& e , const std::string& p){
        User::name = n;
        User::surname = s;
        User::email = e;
        User::auth_token = p;

    }
    std::string User::getsurname(){
        return surname;
    }
    std::string User::getemail(){
        return email;
    }
    std::string User::getName(){
        return name;
    }
    std::string User::getauth_token(){
        return auth_token;
    }
    std::list<Ambiente> * User::getAmbienti(){
        return &ambienti;
    }
    


}  // namespace Server
