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
    bool Queries (std::list<std::string> queries){
        bool correct = true;
         try {

            sql::Driver* driver = get_driver_instance();
            std::unique_ptr<sql::Connection> con(driver->connect(HOST, USER, PASS));
            con->setSchema(DB);
            std::unique_ptr<sql::Statement> stmt(con->createStatement());
            std::list<std::string>::iterator query_iterator;
            for(query_iterator = queries.begin(); query_iterator != queries.end();query_iterator ++ ){
                if (correct == true){
                    std::cout << *query_iterator << std::endl;
                    stmt->execute(*query_iterator);
                }
                else{
                    stmt->execute("rollback");
                    break;
                }
                
            }
            con.reset(nullptr);
            stmt.reset(nullptr);
            } catch (sql::SQLException &e ) {
                correct = false;
                std::cout << "# ERR: SQLException in " << __FILE__;
                std::cout << "(" << __FUNCTION__ << ") on line " << __LINE__ << std::endl;
                std::cout << "# ERR: " << e.what();
                std::cout << " (MySQL error code: " << e.getErrorCode();
                std::cout << ", SQLState: " << e.getSQLState() << " )" << std::endl;
                

            }
        

        return correct;

    }
    int Query( std::string query){
           

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
                

                return 1;
            }

        return 0;
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

