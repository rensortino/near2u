#ifndef SERVER_ATTUATORE_H
#define SERVER_ATTUATORE_H
#include "Dispositivo.hpp"

#include <string>
#include <list>

class Attuatore : public Dispositivo{

    public:

        Attuatore(int code, std::string& type, std::string& name, std::list<std::string> * commands);

        std::list<std::string> * getComandi();

        device_type get_device_type();

        bool controllaComando(std::string& comando);



    private:

        std::list<std::string> comandi;

};

#endif