
#include "User.hpp"


    User::User(const std::string& n, const std::string& s ,const std::string& e , const std::string& p, const std::string& pa){
        name = n;
        surname = s;
        email = e;
        auth_token = p;
        password = pa;

    }
    std::string& User::getsurname(){
        return surname;
    }
    std::string& User::getPassword(){
        return password;
    }
    std::string& User::getemail(){
        return email;
    }
    std::string& User::getName(){
        return name;
    }
    std::string& User::getauth_token(){
        return auth_token;
    }
    std::list<Ambiente *> * User::getAmbienti(){
        return &ambienti;
    }
    Ambiente * User::getAmbiente(std::string& cod_Ambiente){
        for(Ambiente * ambiente : ambienti){
            if(ambiente->getcodAmbiente().compare(cod_Ambiente) == 0){
                return ambiente;
        }
        }

        return nullptr;

    }
    bool User::getAdmin(){
        return admin;
    }
    void User::setAdmin(bool role){
        admin = role;
    }

    void User::addAmbiente(std::string& nome, std::string& codice ){
        Ambiente * ambiente = new Ambiente(nome,codice);
        ambienti.push_back(ambiente);
    }
    
    std::list<Dispositivo *> * User::getDispositivi(std::string& cod_ambiente){
        Ambiente * ambiente;
        ambiente = User::getAmbiente(cod_ambiente);
        if(ambiente != nullptr){
            return ambiente->getDispositivi();
        }
        return nullptr; 
    }
    void User::deleteDispositivo(std::string& cod_ambiente, int code){
        Ambiente * ambiente;
        ambiente = User::getAmbiente(cod_ambiente);
        ambiente ->deleteDispositivo(code);
    }
    
    void User::addDispositivo(std::string& cod_ambiente,int code, std::string& nome, std::string& tipo, std::list<std::string> * commands){
        Ambiente * ambiente = User::getAmbiente(cod_ambiente);
        ambiente->addDispositivo(code,nome,tipo,commands);
        
    }

    bool User::inviaComando(std::string& cod_ambiente, int code_attuatore, std::string& comando){

        Ambiente * ambiente = getAmbiente(cod_ambiente);

        if(ambiente == nullptr){
            return false;
        }
        if(ambiente->inviaComando(code_attuatore,comando) == false){
            return false;
        }
        return true;

    }

    bool User::eliminaAmbiente(std::string& cod_ambiente){
        Ambiente * ambiente;
        ambiente= getAmbiente(cod_ambiente);
        if(ambiente == nullptr){
            return false;
        }
        ambienti.remove(ambiente);
        delete ambiente;
        return true;

    }

    User::~User(){
        
        for(Ambiente * ambiente : ambienti){
           delete ambiente;
        }
        ambienti.clear();

    }
    

