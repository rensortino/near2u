package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func socketReceive(conn net.Conn) []byte {
	var buf bytes.Buffer
	responseSize, err := io.Copy(&buf, conn)
	fmt.Println("total size:", responseSize)
	errorCheck(err, "Error receiving data")
	log.Printf("Receive: %v\n", buf.String())
	return buf.Bytes()
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

