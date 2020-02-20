package utils

import (
	"encoding/json"
	"log"
	"regexp"
)

// TODO Remove nested struct
type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Defines the User data to maintain during a session
type User struct {
	registerData
	Token string
}

// Represents the User data to parse into JSON
type registerData struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(rx chan []byte, name string, surname string, email string, password string)  {

	newUser := registerData {
		name,
		surname,
		email,
		password,
	}

	req := struct {
		Function string 		`json:"function"`
		NewUser registerData    `json:"data"`
	} {
		"register",
		newUser,
	}

	jsonReq, err := json.Marshal(req)
	check(err, "Marshalling Error")
	SocketCommunicate(jsonReq, rx)
}

func Login(responseMsg chan string, token chan string, email string, password string) {

	login := loginData{
		email,
		password,
	}

	req := struct {
		Function string `json:"function"`
		Login  loginData     `json:"data"`
	} {
		"login",
		login,
	}

	jsonReq, err := json.Marshal(req)
	check(err, "Marshalling Error")
	rx := make(chan []byte)
	go SocketCommunicate(jsonReq, rx)
	// TODO Change test string
	//rx <- []byte("0CC0FA6935783505506B0E3B81A566E1B9A7FEBA") // Test SHA1 string
	var res = struct {
		Status string `json:"Status"`
		Message string `json:"message"`
	}{}
	json.Unmarshal(<- rx, &res)
	res.Status = "Succesfull"
	log.Println(res.Status)
	if (res.Status == "Succesfull"){

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
}