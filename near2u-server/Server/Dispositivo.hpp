#ifndef SERVER_DISPOSITIVO_H
#define SERVER_DISPOSITIVO_H
#include<string>

class Dispositivo{

    private:
        int codice;
        std::string tipo;
        std::string nome;



    public:

    Dispositivo(std::string& tipo, std::string& nome, int codice);
    virtual ~Dispositivo();

    int getCodice();

    std::string& getTipo();

    std::string& getNome();
};

#endif