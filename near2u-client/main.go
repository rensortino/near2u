package main

import (
	"./gui"
	"./utils"
	qt "github.com/therecipe/qt/widgets"
)

// TODO implement dynamic binding


func main() {

	gui.InitWindow()

	utils.ClientInstance = new(utils.Client)

	

	qt.QApplication_Exec()

}


