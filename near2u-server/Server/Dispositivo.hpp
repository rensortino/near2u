#ifndef SERVER_DISPOSITIVO_H
#define SERVER_DISPOSITIVO_H
#include<string>
enum device_type {sensore, attuatore};
class Dispositivo{

    private:
        int codice;
        std::string tipo;
        std::string nome;



    public:


    Dispositivo(std::string& tipo, std::string& nome, int codice);
    virtual ~Dispositivo();

    virtual  device_type get_device_type();

    int getCodice();

    std::string& getTipo();

    std::string& getNome();
};

#endif