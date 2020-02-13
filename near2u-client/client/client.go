package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	StopCharacter = "\r\n\r\n"
)

// User represents authentication data
type User struct {
	Name string
	Surname string
	Email    string
	Password string
}

// TODO Find alternative to global variables
var (
	conn net.Conn
	function string
	data string
	auth string
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func GetConnection() net.Conn {
	return conn
}

func SetConnection(newConnection net.Conn) {
	conn = newConnection
}

// SocketConnect binds to a remote socket
func SocketConnect(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	SetConnection(conn)

	check(err)
}

// socketSend sends string on socket connection
func socketSend(conn net.Conn, function string, data string, auth string) string {

	msg := fmt.Sprintf("function:%s data:%s auth:%s%s", function, data, auth, StopCharacter)
	conn.Write([]byte(msg))
	log.Printf("Sending: %s", msg)

	//TODO: segnalare errori nella comunicazione al chiamante

	return "ok"

}

func socketReceive(conn net.Conn, rx chan string) {
	buff := make([]byte, 8192) // Buffered reads from socket

	for {
		n, err := conn.Read(buff)
		if err != nil && err.Error() == "EOF" {
			log.Println("EOF Reached, breaking loop")
			rx <- "EOF"
			break
		}
		check(err)
		log.Printf("Receive: %s\n", buff[:n])
		rx  <- string(buff[:n])
		log.Printf("Channel content: %s\n", <- rx)
	}
}

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

	socketSend(conn, function, data, auth)
	socketReceive(conn, rx)
}

func Login(conn net.Conn, email string, password string) bool {

	userLogin := &User{}
	userLogin.Email = email
	userLogin.Password = password

	req, _ := json.Marshal(userLogin)

	function = "login"
	data = string(req)

	socketSend(conn, function, data, auth)
	rx := make(chan string)
	socketReceive(conn, rx)
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

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func GetEnvList(conn net.Conn) []string {
/*
	function = "getEnvList"
	data = ""
	auth = "a"

	socketSend(conn, function, data, auth)
	SocketReceive(conn)
*/
	//return strings.Split(<- Socket, ";")
	// TODO Delete test string
	var test = make([]string, 10)
	for i := 0; i < 10; i++ {
		test[i] = "Test " + strconv.Itoa(i)
	}
	return test
}

func selectEnv(conn net.Conn, rx chan string, envID string) {

	function = "selectEnvironment"
	data = envID

	socketSend(conn, function, data, auth)
	socketReceive(conn, rx)

	

}