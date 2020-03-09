#include "SHA_CRYPTO.hpp"


    std::string SHA_Crypto( std::string stringa){

       
        unsigned char hash[SHA_DIGEST_LENGTH];
        char digest[41];
        SHA_CTX sha;
        SHA1_Init(&sha);
        SHA1_Update(&sha, stringa.c_str(), stringa.size());
        SHA1_Final(hash, &sha);
        int i = 0;
        for(i = 0; i < SHA_DIGEST_LENGTH; i++)
        {
            sprintf(digest + (i * 2), "%02x", hash[i]);
        }
        digest[40] = 0;
        std::string result(digest);
        return result;


    }
