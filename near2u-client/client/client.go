package client


import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	message       = "Ping"
	StopCharacter = "\r\n\r\n"
)

// User represents authentication data
type User struct {
	Name string
	Surname string
	Email    string
	Password string
}

var (
	Conn net.Conn = nil
)

// SocketConnect binds to a remote socket
func SocketConnect(ip string, port int) net.Conn {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	return conn
}

// SocketCommunicate sends string on socket connection
func SocketCommunicate(conn net.Conn, function string, data string, auth string) {

	msg := fmt.Sprintf("function:%s|data:%s|auth:%s", function, data, auth)
	conn.Write([]byte(msg))
	conn.Write([]byte(StopCharacter))
	log.Printf("Sending: %s", msg)

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Receive: %s", buff[:n])

	//TODO: segnalare errori nella comunicazione al chiamante

}

func Register(conn net.Conn, name string, surname string, email string, password string) bool {

	newUser := &User{}
	newUser.Name = name
	newUser.Surname = surname
	newUser.Email = email
	newUser.Password = password

	// Converts the structure into JSON format
	req, _ := json.Marshal(newUser)

	function := "register"
	data := string(req)
	// TODO substitute with token
	auth := "a"

	SocketCommunicate(conn, function, data, auth)

	return true
}

func Login(conn net.Conn, email string, password string) bool {
	return true
}