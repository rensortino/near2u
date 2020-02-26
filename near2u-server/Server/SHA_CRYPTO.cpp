#include "SHA_CRYPTO.hpp"


    std::string SHA_Crypto( std::string stringa){

        char aux [stringa.length() + 1];
        

        strcpy(aux,stringa.c_str());


        char token[strlen(aux)];

        SHA1((unsigned char *)aux, strlen(aux), (unsigned char *)token);
        char  result [strlen(aux)];
        for(int i = 0; i < strlen(aux) ; i++)
            sprintf(&result[i], "%02x", (unsigned int)token[i]);
        
        std::string token_string (result);
        return token_string;


    }
