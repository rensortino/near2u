package main

import (
	"./gui"
	"./utils"
	qt "github.com/therecipe/qt/widgets"
)

// TODO implement dynamic binding
var (
	ip   = "127.0.0.1"
	port = 3333
)

func main() {

	gui.InitWindow()

	utils.ClientInstance.SocketConnect(ip, port)
	defer utils.ClientInstance.GetConnection().Close()

	qt.QApplication_Exec()

}


