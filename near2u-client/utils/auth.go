package utils

import (
	"encoding/json"
	"log"
	"net"
	"regexp"
)

type LoginData struct {
	Email    string `json:email`
	Password string `json:password`
}

// User represents authentication data
type User struct {
	Name string `json:name`
	Surname string `json:surname`
	LoginData `json:"login"`
}

type RegisterRequest struct {
	RequestParams
	User
}

type LoginRequest struct {
	RequestParams
	LoginData
}

func Register(rx chan []byte, name string, surname string, email string, password string)  {

	loginData := LoginData {
		email,
		password,
	}

	newUser := User{
		name,
		surname,
		loginData,
	}

	params := RequestParams {
		"register",
		"",
	}

	req := RegisterRequest {
		params,
		newUser,
	}

	jsonReq, err := json.Marshal(req)
	check(err, "Marshalling Error")
	SocketSend(conn, jsonReq)
	SocketReceive(conn, rx)
}

func Login(conn net.Conn, email string, password string) bool {

	login := LoginData {
		email,
		password,
	}

	params := RequestParams {
		"login",
		"",
	}

	loginReq := LoginRequest {
		params,
		login,
	}

	jsonReq, err := json.Marshal(loginReq)
	check(err, "Marshalling Error")
	SocketSend(conn, jsonReq)
	rx := make(chan []byte)
	SocketReceive(conn, rx)
	// TODO Change test string
	rx <- []byte("0CC0FA6935783505506B0E3B81A566E1B9A7FEBA") // Test SHA1 string
	var token string
	token = string(<- rx)

	// Checks if token matches the SHA1 format
	isHash, _ := regexp.MatchString("\b[0-9a-f]{5,40}\b", token)

	if isHash {
		auth := token
		log.Println("Token: " + auth)
		return true
	}
	return false
}