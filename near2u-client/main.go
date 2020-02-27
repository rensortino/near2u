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

	selEnvBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getSelectEnvWidget())
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	newEnvBtn := qt.NewQPushButton2("Create Environment", nil)

	newEnvBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getNewEnvWidget())
	})
	layout.AddWidget(newEnvBtn, 0, 0)

	confEnvBtn := qt.NewQPushButton2("Configure Environment", nil)

	confEnvBtn.ConnectClicked(func (checked bool) {
		//changeWindow(widget, getConfEnvWidget())
	})
	layout.AddWidget(confEnvBtn, 0, 0)

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

		responseMsg := make(chan string)
		token := make(chan string)
		go utils.Login(responseMsg, token, email.Text(), password.Text())
		clientInstance.LoggedUser = <- token
		res := <- responseMsg
		if res == "Succesfull" {
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		} else {
			qt.QMessageBox_Information(nil, "Error", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
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
		rx := make(chan string)
		go utils.Register(rx, name.Text(), surname.Text(), email.Text(), password.Text())

		res := <- rx
		if res == "Succesfull"{
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		} else {
			qt.QMessageBox_Information(nil, "Error", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
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

func getSelectEnvWidget() * qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	/*
	envList := qt.NewQComboBox(nil)
	envList.AddItems(clientInstance.GetEnvList())
	envList.SetEditable(false)
	layout.AddWidget(envList, 0, 0)
	*/

	envName := qt.NewQLineEdit(nil)
	layout.AddWidget(envName, 0, 0)

	selEnvBtn := qt.NewQPushButton2("Select", nil)
	selEnvBtn.ConnectClicked(func (checked bool) {
		log.Printf("Selected: %s", envName.Text())
		topicCh := make(chan string) // Stores the topic returned from the server
		errCh := make(chan string) // Stores the error message, in case of failed request
		go clientInstance.SelectEnv(envName.Text(), topicCh, errCh)
		select {
			case topic := <- topicCh :
				changeWindow(widget, getRTDataWidget(topic))
			case error := <- errCh:
				qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getNewEnvWidget() * qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	envName := qt.NewQLineEdit(nil)
	layout.AddWidget(envName, 0, 0)

	newEnvBtn := qt.NewQPushButton2("Create Environment", nil)
	newEnvBtn.ConnectClicked(func (checked bool) {
		envCh := make(chan * client.Environment) // Stores the environment returned from the server
		errCh := make(chan string) // Stores the error message, in case of failed request
		go clientInstance.CreateEnv(envName.Text(), envCh, errCh)
		select {
		case newEnv := <- envCh :
			changeWindow(widget, getAddSensorsWidget(newEnv))
		case error := <- errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(newEnvBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getRTDataWidget(topic string) * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	dataList := qt.NewQListWidget(nil)

	rtCh := make(chan map[string]interface{}) // Channel for real time data
	quit := make(chan bool)

	clientInstance.GetSensorData(topic, rtCh)

	// Used a goroutine to run the subscribe callback in parallel and update the data
	go func() {
		for {
			select {
			case receivedData := <-rtCh:
				for id, sensor := range receivedData {
					name := sensor.(map[string]interface{})["Name"].(string)
					measurement := sensor.(map[string]interface{})["Measurement"].(float64)
					dataList.AddItem(fmt.Sprintf("%s\t%s\t%f\n",id, name, measurement))
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

func getAddSensorsWidget(newEnv * client.Environment) * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	code := qt.NewQLineEdit(nil)
	code.SetPlaceholderText("Sensor Code")
	layout.AddWidget(code, 0, 0)

	name := qt.NewQLineEdit(nil)
	name.SetPlaceholderText("Sensor Name")
	layout.AddWidget(name, 0, 0)

	kind := qt.NewQLineEdit(nil)
	kind.SetPlaceholderText("Sensor Kind")
	layout.AddWidget(kind, 0, 0)

	addBtn := qt.NewQPushButton2("Add", nil)
	addBtn.ConnectClicked(func (checked bool) {
		sensorCh := make(chan interface{})
		errCh := make(chan string) // Stores the error message, in case of failed request
		go clientInstance.AddSensor(code.Text(), name.Text(), kind.Text(), newEnv, sensorCh, errCh)
		select {
		case newSensor := <- sensorCh :
			successString := fmt.Sprintf("Sensor: %s Added", newSensor.(* client.Sensor).Name)
			qt.QMessageBox_Information(nil, "OK", successString, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		case error := <- errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(addBtn, 0 , 0)

	doneBtn := qt.NewQPushButton2("Done", nil)
	doneBtn.ConnectClicked(func (checked bool) {
		resCh := make(chan string)
		errCh := make(chan string)
		go clientInstance.Done(newEnv, resCh, errCh)
		select {
			case res := <- resCh:
				qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			case err := <- errCh:
				qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(doneBtn, 0 , 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func (checked bool) {
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