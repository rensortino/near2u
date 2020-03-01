#ifndef SERVER_SENSORE_H
#define SERVER_SENSORE_H

#include <string>
#include <vector>
#include <list>
#include <iostream>
#include <assert.h>


class Sensore
{
public:
	Sensore(int cod, std::string& nome, std::string& tipo);
	int getCodSensore();
	std::string& getName();
	std::string& getType();
private:
	int codSensore;
	std::string type;
	std::string name;

};

#endif
