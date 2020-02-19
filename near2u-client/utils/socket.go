package utils

import (
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	ip   = "127.0.0.1"
	port = 3333
)

var ClientInstance * Client

// TODO Convert in Session interface
type Client struct {
	ClientID string
	conn net.Conn
	Token string
}

type RequestParams struct {
	Function string `json:"function"`
	Auth string `json:"auth"`
}

func check(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Fatalln(err)
	}
}

// TODO make methods private

func (c * Client) GetConnection() net.Conn {
	return c.conn
}

func (c * Client) SetConnection(newConnection net.Conn) {
	c.conn = newConnection
}

// SocketConnect binds to a remote socket
func (c * Client) SocketConnect(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	c.SetConnection(conn)

	check(err, "Connection Failed")
}

// socketSend sends JSON on socket connection, accepts marshaled JSON
func (c * Client) SocketSend(jsonReq []byte) {
	c.SocketConnect(ip, port)
	
	_, err := c.conn.Write(jsonReq)
	log.Printf("Sending: %s", string(jsonReq))

	check(err, "Couldn't send data")
}

func (c * Client) SocketReceive(rx chan []byte) {
	buff := make([]byte, 8192) // Buffered reads from socket
	defer c.GetConnection().Close()
	for {
		n, err := c.conn.Read(buff)
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