#include <string>
#include <vector>
#include <list>
#include <iostream>
#include <assert.h>

#include "Sensore.h"


Sensore::Sensore(int code, std::string& name,std::string& type) : Dispositivo(type,name,code){
}

device_type Sensore::get_device_type(){
    return sensore;

}