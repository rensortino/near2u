package utils

import (
	"log"
	"net"
	"strconv"
	"strings"
)

// TODO Make dynamic
var (
	ip   = "127.0.0.1"
	port = 3333
)

func errorCheck(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Fatalln(err)
	}
}

// TODO make methods private

// SocketConnect binds to a remote socket
func socketConnect(ip string, port int) net.Conn{
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	errorCheck(err, "Connection Failed")

	return conn
}

// Sends JSON on socket connection, accepts marshaled JSON
func socketSend(conn net.Conn, jsonReq []byte) {
	
	_, err := conn.Write(jsonReq)
	log.Printf("Sending: %s", string(jsonReq))

	errorCheck(err, "Couldn't send data")
}

func socketReceive(conn net.Conn, rx chan []byte) {
	buff := make([]byte, 8192) // Buffered reads from socket
	for {
		n, err := conn.Read(buff)
		if err != nil && err.Error() == "EOF" {
			log.Println("EOF Reached, breaking loop")
			rx <- []byte("EOF")
			break
		}
		errorCheck(err, "Error receiving data")
		log.Printf("Receive: %s\n", buff[:n])
		rx  <- buff[:n]
	}
}

func SocketCommunicate(jsonReq []byte, rx chan []byte) {

	conn := socketConnect(ip, port)
	defer conn.Close()

	socketSend(conn, jsonReq)
	socketReceive(conn, rx)
}