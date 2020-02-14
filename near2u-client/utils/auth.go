package utils

import (
	"encoding/json"
	"log"
	"net"
	"regexp"
)

// User represents authentication data
type User struct {
	Name string
	Surname string
	Email    string
	Password string
}

var (
	function string
	data string
	auth string
)

func Register(rx chan string, name string, surname string, email string, password string)  {

	newUser := &User{}
	newUser.Name = name
	newUser.Surname = surname
	newUser.Email = email
	newUser.Password = password

	// Converts the structure into JSON format
	req, _ := json.Marshal(newUser)

	function = "register"
	data = string(req)

	SocketSend(conn, function, data, auth)
	SocketReceive(conn, rx)
}

func Login(conn net.Conn, email string, password string) bool {

	userLogin := &User{}
	userLogin.Email = email
	userLogin.Password = password

	req, _ := json.Marshal(userLogin)

	function = "login"
	data = string(req)

	SocketSend(conn, function, data, auth)
	rx := make(chan string)
	SocketReceive(conn, rx)
	// TODO Change test string
	rx <- "0CC0FA6935783505506B0E3B81A566E1B9A7FEBA" // Test SHA1 string
	var token string
	token = <- rx

	// Checks if token matches the SHA1 format
	isHash, _ := regexp.MatchString("\b[0-9a-f]{5,40}\b", token)

	if isHash {
		auth = token
		log.Println("Token: " + auth)
		return true
	}
	return false
}