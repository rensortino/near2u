package utils

import (
	"encoding/json"
	"log"
	"regexp"
)

// TODO Remove nested struct
type loginData struct {
	Email    string `json:email`
	Password string `json:password`
}

type user struct {
	Name      string `json:name`
	Surname   string `json:surname`
	LoginData loginData `json:"login"`
}

func (c * Client) Register(rx chan []byte, name string, surname string, email string, password string)  {

	login := loginData{
		email,
		password,
	}

	newUser := user{
		name,
		surname,
		login,
	}

	params := RequestParams {
		"register",
		"",
	}

	req := struct {
		Params  RequestParams `json:"params"`
		NewUser user          `json:"new_user"`
	} {
		params,
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

	params := RequestParams {
		"login",
		"",
	}

	req := struct {
		Params RequestParams `json:"params"`
		Login  loginData     `json:"login"`
	} {
		params,
		login,
	}

	jsonReq, err := json.Marshal(req)
	check(err, "Marshalling Error")
	c.SocketSend(jsonReq)
	rx := make(chan []byte)
	c.SocketReceive(rx)
	// TODO Change test string
	rx <- []byte("0CC0FA6935783505506B0E3B81A566E1B9A7FEBA") // Test SHA1 string
	var token string
	token = string(<- rx)

	// Checks if token matches the SHA1 format
	isHash, _ := regexp.MatchString("\b[0-9a-f]{5,40}\b", token)

	if isHash {
		c.Token = token
		log.Println("Token: " + token)
		return true
	}
	return false
}