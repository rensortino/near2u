package utils

import (
	"encoding/json"
	"log"
	//"regexp"
)

// TODO Remove nested struct
type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type user_register struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c * Client) Register(rx chan []byte, name string, surname string, email string, password string)  {

	newUser := user_register{
		name,
		surname,
		email,
		password,
	}

	req := struct {
		Function string `json:"function"`
		NewUser user_register     `json:"data"`
	} {
		"register",
		newUser,
	}

	jsonReq, err := json.Marshal(req)
	check(err, "Marshalling Error")
	c.SocketSend(jsonReq)
	c.SocketReceive(rx)
}

func (c * Client) Login(email string, password string) bool {

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
	c.SocketSend(jsonReq)
	rx := make(chan []byte)
	go c.SocketReceive(rx)
	// TODO Change test string
	//rx <- []byte("0CC0FA6935783505506B0E3B81A566E1B9A7FEBA") // Test SHA1 string
	var res = struct {
		Status string `json:"Status"`
		Message string `json:"message"`
	}{}
	json.Unmarshal(<- rx,&res)
	log.Println(res.Status)
	if (res.Status == "Succesfull"){
		c.Token = res.Message
		log.Println(c.Token)
	} else
		{
			// da mettere nella gui il messaggio di errore 
			log.Println("error")
			c.Token = ""
		}

	// Checks if token matches the SHA1 format
	/*
	isHash, _ := regexp.MatchString("\b[0-9a-f]{5,40}\b", token)

	if isHash {
		c.Token = token
		log.Println("Token: " + token)
		return true
	}
	*/
	return false
}