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

	utils.SocketConnect(ip, port)
	defer utils.GetConnection().Close()

	qt.QApplication_Exec()

}


