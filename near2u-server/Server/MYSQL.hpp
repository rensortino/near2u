#include <string>
#include <cppconn/resultset.h>
#include <jsoncpp/json/json.h>

namespace MYSQL{

int Query( std::string query);
bool Queries (std::list<std::string> queries);
sql::ResultSet* Select_Query(std::string query);

}