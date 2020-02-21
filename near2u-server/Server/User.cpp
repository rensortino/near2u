
#include "User.hpp"


    User::User(const std::string& n, const std::string& s ,const std::string& e , const std::string& p, const std::string& pa){
        name = n;
        surname = s;
        email = e;
        auth_token = p;
        password = pa;

    }
    std::string User::getsurname(){
        return surname;
    }
    std::string User::getPassword(){
        return password;
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
    

