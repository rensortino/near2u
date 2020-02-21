package main

import (
	"./client"
	"./utils"
	"fmt"
	qtcore "github.com/therecipe/qt/core"
	qtgui "github.com/therecipe/qt/gui"
	qt "github.com/therecipe/qt/widgets"
	"log"
	"os"
	"time"
)

var (
	mainWindow * qt.QMainWindow
	mainWidget * qt.QWidget
	clientInstance * client.Client
)

func changeWindow(oldWindow * qt.QWidget, newWindow * qt.QWidget) {
	oldWindow.DestroyQWidget()
	mainWindow.SetCentralWidget(newWindow)
}

func getHomepageWidget() * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)


	loginBtn := qt.NewQPushButton2("Log In", nil)
	loginBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getLoginWidget())
	})
	layout.AddWidget(loginBtn, 0, 0)

	registerBtn := qt.NewQPushButton2("Register", nil)
	registerBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getRegisterWidget())
	})
	layout.AddWidget(registerBtn, 0, 0)

	selEnvBtn := qt.NewQPushButton2("Select Environment", nil)
	/*
	if clientInstance.LoggedUser == "" {
		selEnvBtn.SetDisabled(true)
	}
	 */
	selEnvBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getSelectEnvWidget())
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	return widget
}

func getSelectEnvWidget() * qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	envList := qt.NewQComboBox(nil)
	envList.AddItems(clientInstance.GetEnvList())
	envList.SetEditable(false)
	layout.AddWidget(envList, 0, 0)

	selEnvBtn := qt.NewQPushButton2("Select", nil)
	selEnvBtn.ConnectClicked(func (checked bool) {
		log.Printf("Selected: %s", envList.CurrentText())

		// TODO Change
		clientInstance.SelectEnv("env1")
		changeWindow(widget, getRTDataWidget())
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getRegisterWidget() * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	name := qt.NewQLineEdit(nil)
	name.SetPlaceholderText("Name")
	layout.AddWidget(name, 0, 0)

	surname := qt.NewQLineEdit(nil)
	surname.SetPlaceholderText("Surname")
	layout.AddWidget(surname, 0, 0)

	email := qt.NewQLineEdit(nil)
	email.SetPlaceholderText("Email")

	regex := qtcore.NewQRegExp()
	regex.SetPattern("\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,4}\\b")
	layout.AddWidget(email, 0, 0)
	validator := qtgui.NewQRegExpValidator2(regex, nil)
	email.SetValidator(validator)

	password := qt.NewQLineEdit(nil)
	password.SetPlaceholderText("Password")
	layout.AddWidget(password, 0, 0)

	registerBtn := qt.NewQPushButton2("Register", nil)
	registerBtn.ConnectClicked(func(checked bool) {
		// TODO Email field validation check
		rx := make(chan []byte)
		go utils.Register(rx, name.Text(), surname.Text(), email.Text(), password.Text())
		// TODO Verify server response correctness
		qt.QMessageBox_Information(nil, "OK", string(<- rx), qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		changeWindow(widget, getHomepageWidget())

	})
	layout.AddWidget(registerBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getLoginWidget() * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	email := qt.NewQLineEdit(nil)
	email.SetPlaceholderText("Email")
	layout.AddWidget(email, 0, 0)

	password := qt.NewQLineEdit(nil)
	password.SetPlaceholderText("Password")
	layout.AddWidget(password, 0, 0)

	loginBtn := qt.NewQPushButton2("Log In", nil)
	loginBtn.ConnectClicked(func(checked bool) {
		// TODO Fields validity check
		rx := make(chan string)
		token := make(chan string)
		go utils.Login(rx, token, email.Text(), password.Text())
		// TODO Assign token to instance
		log.Println(<- token)
		qt.QMessageBox_Information(nil, "OK", <- rx, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(loginBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getRTDataWidget() * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	dataList := qt.NewQListWidget(nil)

	// TODO Parametrize topic
	topic := "testtopic"
	rtCh := make(chan map[string]client.Sensor) // Channel for real time data
	quit := make(chan bool)

	clientInstance.GetSensorData(topic, rtCh)

	// Used a goroutine to run the subscribe callback in parallel and update the data
	go func() {
		for {
			select {
			case receivedData := <-rtCh:
				for id, sensor := range receivedData {
					dataList.AddItem(fmt.Sprintf("%s\t%s\t%f\n",id, sensor.Name, sensor.Measurement))
				}
			case <-quit:
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	layout.AddWidget(dataList, 0 , 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
		clientInstance.StopGettingData(topic, rtCh, quit)
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func initWindow() {

	// Creates a new graphic application
	qt.NewQApplication(len(os.Args), os.Args)

	mainWindow = qt.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle("Near2U")
	mainWindow.SetMinimumSize2(400, 400)

	mainWidget = getHomepageWidget()

	mainWindow.SetCentralWidget(mainWidget)
	mainWindow.Show()
}

func main() {

	clientInstance = client.GetClientInstance()
	initWindow()
	qt.QApplication_Exec()

}