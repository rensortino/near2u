package utils

import (
	"log"
)

func Register(rx chan interface{}, name string, surname string, email string, password string)  {

	newUser := struct {
		Name string `json:"name"`
		Surname string `json:"surname"`
		Email string `json:"email"`
		Password string `json:"password"`
	} {
		name,
		surname,
		email,
		password,
	}

	SocketCommunicate("register", "", newUser, rx)
}

func Login(responseMsg chan string, token chan string, email string, password string) {

	data := struct {
		Email string
		Password string
	} {
		email,
		password,
	}

	rx := make(chan interface{})
	go SocketCommunicate("login", "", data, rx)
	// TODO Change test string
	//rx <- []byte("0CC0FA6935783505506B0E3B81A566E1B9A7FEBA") // Test SHA1 string

	res := (<- rx).(map[string]interface{}) // res has type map[string]interface{}

	log.Println(res["status"])
	if res["status"] == "Succesfull" {

		token <- res["message"].(string)
		responseMsg <- "User Authenticated"
	}else {
		token <- "NULL"
		responseMsg <- res["message"].(string)
	}

		/*

		// Checks if token matches the SHA1 format
		isHash, _ := regexp.MatchString("\b[0-9a-f]{5,40}\b", res.Message)

		if isHash {
			token <- res.Message
			log.Println("Token sent from server: " + res.Message)
			responseMsg <- "User Authenticated"
		} else {
			token <- "NULL"
			responseMsg <- "Token not valid"
		}
	} else
	{
		log.Println(res)
		token <- "NULL"
		responseMsg <- "Error"
	}
	*/

}