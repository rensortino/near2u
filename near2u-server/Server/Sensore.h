#ifndef SERVER_SENSORE_H
#define SERVER_SENSORE_H
#include "Dispositivo.hpp"
#include <string>
#include <vector>
#include <list>
#include <iostream>
#include <assert.h>


class Sensore : public Dispositivo
{
public:
	Sensore(int cod, std::string& name, std::string& type);

	device_type get_device_type();

};

#endif
