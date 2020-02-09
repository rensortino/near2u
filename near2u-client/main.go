package main

import (
	"./client"
	"./gui"
	qt "github.com/therecipe/qt/widgets"
)

// TODO implement dynamic binding
var (
	ip   = "127.0.0.1"
	port = 3333
)

func main() {

	gui.InitWindow()

	client.Conn = client.SocketConnect(ip, port)
	defer client.Conn.Close()

	qt.QApplication_Exec()

}


