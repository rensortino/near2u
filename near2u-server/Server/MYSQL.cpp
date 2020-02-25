#include "MYSQL.hpp"
#include <stdlib.h>
#include <iostream>
#include <sstream>
#include <stdexcept>

#include "mysql_connection.h"

#include <cppconn/driver.h>
#include <cppconn/exception.h>
#include <cppconn/resultset.h>
#include <cppconn/statement.h>
#include <cppconn/prepared_statement.h>
#include "User.hpp"
#include <jsoncpp/json/json.h>

#define HOST "localhost"
#define USER "admin"
#define PASS "admin"
#define DB "apl_project"

namespace MYSQL{
    Json::Value insert( std::string query){
            Json::Value response;
            response["status"] = "";
		    response["error"] = "";

        try {

            sql::Driver* driver = get_driver_instance();
            std::unique_ptr<sql::Connection> con(driver->connect(HOST, USER, PASS));
            con->setSchema(DB);
            std::unique_ptr<sql::Statement> stmt(con->createStatement());
            stmt->execute(query);
            con.reset(nullptr);
            stmt.reset(nullptr);
            } catch (sql::SQLException &e) {
        

                std::cout << "# ERR: SQLException in " << __FILE__;
                std::cout << "(" << __FUNCTION__ << ") on line " << __LINE__ << std::endl;
                std::cout << "# ERR: " << e.what();
                std::cout << " (MySQL error code: " << e.getErrorCode();
                std::cout << ", SQLState: " << e.getSQLState() << " )" << std::endl;

                response["status"] = "Error";
                response["error"] = e.what();
                

                return response;
            }

        response["status"] = "Succesfull";
        response["error"] = "";
        return response;
    }

    sql::ResultSet* Select_Query(std::string query){
            sql::ResultSet  *res;
        try {
            sql::Driver* driver = get_driver_instance();
            std::unique_ptr<sql::Connection> con(driver->connect(HOST, USER, PASS));
            con->setSchema(DB);
            std::unique_ptr<sql::Statement> stmt(con->createStatement());
            res = stmt->executeQuery(query);
            con.reset(nullptr);
            stmt.reset(nullptr);
            
            } catch (sql::SQLException &e) {
        

                std::cout << "# ERR: SQLException in " << __FILE__;
                std::cout << "(" << __FUNCTION__ << ") on line " << __LINE__ << std::endl;
                std::cout << "# ERR: " << e.what();
                std::cout << " (MySQL error code: " << e.getErrorCode();
                std::cout << ", SQLState: " << e.getSQLState() << " )" << std::endl;

                return nullptr;
            }

        return res;
    }
}

