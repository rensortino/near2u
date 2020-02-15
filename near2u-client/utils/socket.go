package utils

import (
	"log"
	"net"
	"strconv"
	"strings"
)

var conn net.Conn

type Client struct {
	ClientID string
	Conn net.Conn
	Token string
}

type RequestParams struct {
	Function string `json:function`
	Auth string `json:auth`
}

// TODO Remove (?)
const StopCharacter = "\r\n\r\n"

func check(err error, msg string) {
	if err != nil {
		log.Println(msg)
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

	check(err, "Connection Failed")
}

// socketSend sends JSON on socket connection, accepts marshaled JSON
func SocketSend(conn net.Conn, jsonReq []byte) {

	_, err := conn.Write(jsonReq)
	log.Printf("Sending: %s", string(jsonReq))

	check(err, "Couldn't send data")
}

func SocketReceive(conn net.Conn, rx chan []byte) {
	buff := make([]byte, 8192) // Buffered reads from socket

	for {
		n, err := conn.Read(buff)
		if err != nil && err.Error() == "EOF" {
			log.Println("EOF Reached, breaking loop")
			rx <- []byte("EOF")
			break
		}
		check(err, "Error receiving data")
		log.Printf("Receive: %s\n", buff[:n])
		rx  <- buff[:n]
		log.Printf("Channel content: %s\n", <- rx)
	}
}