#include "SHA_CRYPTO.h"

namespace Server{
    std::string SHA_Crypto( std::string stringa){

        char aux [stringa.length() + 1];
        

        strcpy(aux,stringa.c_str());


        char token[strlen(aux)];

        SHA1((unsigned char *)aux, strlen(aux), (unsigned char *)token);

        char  result [strlen(aux)];
        for(int i = 0; i < 30 ; i++)
            sprintf(&result[i*2], "%02x", (unsigned int)token[i]);
        
        std::string token_string (token);
        return token_string;


    }
}
