#include <string>
#include <cppconn/resultset.h>
#include <jsoncpp/json/json.h>

namespace MYSQL{

Json::Value insert( std::string query);
sql::ResultSet* Select_Query(std::string query);

}