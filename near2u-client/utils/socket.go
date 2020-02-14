package utils

import (
	"log"
	"net"
	"strconv"
	"strings"
)

var conn net.Conn

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

// socketSend sends string on socket connection
func SocketSend(conn net.Conn, function string, data string, auth string) string {

	msg := fmt.Sprintf("function:%s data:%s auth:%s%s", function, data, auth, StopCharacter)
	_, err := conn.Write([]byte(msg))
	log.Printf("Sending: %s", msg)

	check(err, "Couldn't send data")
}

func SocketReceive(conn net.Conn, rx chan string) {
	buff := make([]byte, 8192) // Buffered reads from socket

	for {
		n, err := conn.Read(buff)
		if err != nil && err.Error() == "EOF" {
			log.Println("EOF Reached, breaking loop")
			rx <- "EOF"
			break
		}
		check(err, "Error receiving data")
		log.Printf("Receive: %s\n", buff[:n])
		rx  <- string(buff[:n])
		log.Printf("Channel content: %s\n", <- rx)
	}
}