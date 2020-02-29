package utils

import (
	"encoding/json"
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

type Request struct {
	Function string      `json:"function"`
	Data     interface{} `json:"data"`
	Auth     string      `json:"auth"`
}

func buildRequest(function, auth string, data interface{}) interface{} {
	request := struct {
		Function string      `json:"function"`
		Data     interface{} `json:"data"`
		Auth     string      `json:"auth"`
	}{
		function,
		data,
		auth,
	}

	return request

}

func errorCheck(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Fatalln(err)
	}
}

// SocketConnect binds to a remote socket
func socketConnect(ip string, port int) net.Conn {
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

// TODO Handle messages bigger than buffer
func socketReceive(conn net.Conn) []byte {
	buff := make([]byte, 8192) // Buffered reads from socket
	for {
		n, err := conn.Read(buff)
		if err != nil && err.Error() == "EOF" {
			log.Println("EOF Reached, breaking loop")
			return []byte("EOF")
		}
		errorCheck(err, "Error receiving data")
		log.Printf("Receive: %s\n", buff[:n])
		return buff[:n]
	}
	return []byte("EOF")
}

// Accepts request parameters, returns JSON as map[string]interface{} on the channel
func SocketCommunicate(function, auth string, data interface{}) map[string]interface{} {

	conn := socketConnect(ip, port)
	defer conn.Close()

	req := buildRequest(function, auth, data)

	jsonReq, _ := json.Marshal(req)

	socketSend(conn, jsonReq)
	res := socketReceive(conn)

	jsonRes := make(map[string]interface{})

	if string(res) != "EOF" {
		json.Unmarshal(res, &jsonRes)
	} else {
		eof := []byte(`{"status":"EOF reached"}`)
		json.Unmarshal(eof, &jsonRes)
	}
	return jsonRes
}
