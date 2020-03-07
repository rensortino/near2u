package main

import (
	"./client"
	"./utils"
	"fmt"
	qtcore "github.com/therecipe/qt/core"
	qtgui "github.com/therecipe/qt/gui"
	qt "github.com/therecipe/qt/widgets"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

var (
	mainWindow     *qt.QMainWindow
	mainWidget     *qt.QWidget
	clientInstance *client.Client
	currentEnv     *client.Environment
)

//TODO fix back buttons

func initWindow() {

	// Creates a new graphic application
	qt.NewQApplication(len(os.Args), os.Args)

	mainWindow = qt.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle("Near2U")
	mainWindow.SetMinimumSize2(800, 600)

	mainWidget = getStartWidget()

	mainWindow.SetCentralWidget(mainWidget)
	mainWindow.Show()
}

func changeWindow(oldWindow *qt.QWidget, newWindow *qt.QWidget) {
	oldWindow.DestroyQWidget()
	mainWindow.SetCentralWidget(newWindow)
}

func getStartWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	loginBtn := qt.NewQPushButton2("Log In", nil)
	loginBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getLoginWidget())
	})
	layout.AddWidget(loginBtn, 0, 0)

	registerBtn := qt.NewQPushButton2("Register", nil)
	registerBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getRegisterWidget())
	})
	layout.AddWidget(registerBtn, 0, 0)

	return widget
}

func getHomepageWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	selEnvBtn := qt.NewQPushButton2("Get Real Time Data", nil)
	selEnvBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getSelectEnvForMQTTWidget())
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	if clientInstance.LoggedUser.IsAdmin {
		newEnvBtn := qt.NewQPushButton2("Create / Delete Environment", nil)
		newEnvBtn.ConnectClicked(func(checked bool) {
			changeWindow(widget, getAddDelEnvWidget())
		})
		layout.AddWidget(newEnvBtn, 0, 0)

		confEnvBtn := qt.NewQPushButton2("Configure Environment", nil)
		confEnvBtn.ConnectClicked(func(checked bool) {
			changeWindow(widget, getConfigureEnvWidget())
		})
		layout.AddWidget(confEnvBtn, 0, 0)
	}

	sendCmdBtn := qt.NewQPushButton2("Send Command", nil)
	sendCmdBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getSendCommandWidget())
	})
	layout.AddWidget(sendCmdBtn, 0, 0)

	getHistoryBtn := qt.NewQPushButton2("Visualize History Data", nil)
	getHistoryBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHistoryDataWidget())
	})
	layout.AddWidget(getHistoryBtn, 0, 0)

	logoutBtn := qt.NewQPushButton2("Log Out", nil)
	logoutBtn.ConnectClicked(func(checked bool) {
		resCh := make(chan string)
		errCh := make(chan string)
		go clientInstance.LoggedUser.Logout(resCh, errCh)
		select {
		case res := <- resCh:
			clientInstance.LoggedUser = &utils.User{}
			currentEnv = client.NewEnvironment()
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getStartWidget())
		case err := <- errCh:
			qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}

	})
	layout.AddWidget(logoutBtn, 0, 0)

	return widget
}

func getLoginWidget() *qt.QWidget {

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

		resCh := make(chan string)
		errCh := make(chan string)
		userCh := make(chan * utils.User)
		go utils.Login(resCh, errCh, userCh, email.Text(), password.Text())
		select {
		case res := <- resCh:
			clientInstance.LoggedUser = <- userCh
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getHomepageWidget())
		case err := <- errCh:
			qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(loginBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getRegisterWidget() *qt.QWidget {

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

		res := <-rx
		if res == "User Registered" {
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		} else {
			qt.QMessageBox_Information(nil, "Error", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
		changeWindow(widget, getStartWidget())

	})
	layout.AddWidget(registerBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getSelectEnvForMQTTWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	envListCB := qt.NewQComboBox(nil)
	envListCB.SetEditable(false)

	getEnvList(envListCB)
	layout.AddWidget(envListCB, 0, 0)

	selEnvBtn := qt.NewQPushButton2("Select Environment", nil)
	selEnvBtn.ConnectClicked(func(checked bool) {
		topicCh := make(chan string) // Stores the topic returned from the server
		uriCh := make(chan string)
		errCh := make(chan string) // Stores the error message, in case of failed request
		go clientInstance.GetTopicAndUri(envListCB.CurrentText(), topicCh, uriCh, errCh)
		select {
		case topic := <-topicCh:
			uri := <-uriCh
			changeWindow(widget, getRTDataWidget(topic, uri))
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getConfigureEnvWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	envListCB := qt.NewQComboBox(nil)
	envListCB.SetEditable(false)

	getEnvList(envListCB)
	layout.AddWidget(envListCB, 0, 0)

	addBtn := qt.NewQPushButton2("Add Devices", nil)
	addBtn.ConnectClicked(func(checked bool) {
		client.SetCurrentEnv(currentEnv, envListCB.CurrentText())
		fmt.Printf("\nCURRENTENV: %v", currentEnv)
		resCh := make(chan string)
		errCh := make(chan string)
		go currentEnv.GetDevicesList(resCh, errCh)
		select {
		case res := <- resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getSelectDeviceTypeWidget())
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(addBtn, 0, 0)

	delBtn := qt.NewQPushButton2("Delete Devices", nil)
	delBtn.ConnectClicked(func(checked bool) {
		client.SetCurrentEnv(currentEnv, envListCB.CurrentText())
		resCh := make(chan string)
		errCh := make(chan string)
		go currentEnv.GetDevicesList(resCh, errCh)
		select {
		case res := <- resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getDeleteSensorsWidget())
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(delBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getAddDelEnvWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	envName := qt.NewQLineEdit(nil)
	layout.AddWidget(envName, 0, 0)

	newEnvBtn := qt.NewQPushButton2("Create Environment", nil)
	newEnvBtn.ConnectClicked(func(checked bool) {
		resCh := make(chan string) // Stores the environment returned from the server
		errCh := make(chan string)              // Stores the error message, in case of failed request
		go clientInstance.CreateEnv(envName.Text(), currentEnv, resCh, errCh)
		select {
		case res := <- resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getSelectDeviceTypeWidget())
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(newEnvBtn, 0, 0)

	delEnvBtn := qt.NewQPushButton2("Delete Environment", nil)
	delEnvBtn.ConnectClicked(func(checked bool) {
		resCh := make(chan string) // Stores the environment returned from the server
		errCh := make(chan string)              // Stores the error message, in case of failed request
		go clientInstance.DeleteEnv(envName.Text(), currentEnv, resCh, errCh)
		select {
		case res := <- resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getHomepageWidget())
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(delEnvBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getRTDataWidget(topic, uriString string) *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	dataList := qt.NewQListWidget(nil)
	dataList.AddItem(fmt.Sprintf("Code\tName\tMeasurement\n"))

	rtCh := make(chan interface{}) // Channel for real time data
	quit := make(chan bool)
	startCh := make(chan bool)

	go clientInstance.GetSensorData(topic, rtCh, startCh)

	// Used a goroutine to run the subscribe callback in parallel and update the data
	go func() {
		uri, err := url.Parse("tcp://" + uriString)
		if err != nil {
			qt.QMessageBox_Information(nil, "Error", "Error parsing MQTT URL", qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			return
		}
		// Handles MQTT server connection
			clientInstance.MQTTClient = utils.MQTTConnect(clientInstance.ID, uri)
			startCh <- true
		// Handles received data from subscribed topic
		for {
			select {
			case receivedData, ok := <-rtCh:
				if ok {
					dataList.AddItem(fmt.Sprintf("%d\t%s\t%f\n", receivedData.(*client.Sensor).Code,
						receivedData.(*client.Sensor).Name, receivedData.(*client.Sensor).Measurement))
				}
			case <-quit:
				return
			}
		}
	}()
	layout.AddWidget(dataList, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		clientInstance.StopGettingData(topic, rtCh, quit)
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getSelectDeviceTypeWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	addSensorBtn := qt.NewQPushButton2("Add Sensor", nil)
	addSensorBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getAddSensorActuatorWidget(false))
	})
	layout.AddWidget(addSensorBtn, 0, 0)

	addActuatorBtn := qt.NewQPushButton2("Add Actuator", nil)
	addActuatorBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getAddSensorActuatorWidget(true))
	})
	layout.AddWidget(addActuatorBtn, 0, 0)

	return widget
}

func getAddSensorActuatorWidget(addActuator bool) *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	devicesCB := qt.NewQComboBox(nil)
	devicesCB.SetEditable(false)

	getDeviceList(devicesCB)
	layout.AddWidget(devicesCB, 0, 0)

	code := qt.NewQLineEdit(nil)
	code.SetPlaceholderText("Code")
	layout.AddWidget(code, 0, 0)

	name := qt.NewQLineEdit(nil)
	name.SetPlaceholderText("Name")
	layout.AddWidget(name, 0, 0)

	kind := qt.NewQLineEdit(nil)
	kind.SetPlaceholderText("Kind")
	layout.AddWidget(kind, 0, 0)

	if addActuator {
		commands := make([]string, 0)

		cmd := qt.NewQLineEdit(nil)
		cmd.SetPlaceholderText("Command")
		layout.AddWidget(cmd, 0, 0)

		addCmdBtn := qt.NewQPushButton2("Add Command", nil)
		addCmdBtn.ConnectClicked(func(checked bool) {
			commands = append(commands, cmd.Text())
			qt.QMessageBox_Information(nil, "OK", "Command Added", qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		})
		layout.AddWidget(addCmdBtn, 0, 0)

		clearBtn := qt.NewQPushButton2("Clear Commands", nil)
		clearBtn.ConnectClicked(func(checked bool) {
			commands = make([]string, 0)
		})
		layout.AddWidget(clearBtn, 0, 0)

		addActBtn := qt.NewQPushButton2("Add Actuator", nil)
		addActBtn.ConnectClicked(func(checked bool) {

			if len(commands) == 0 {
				qt.QMessageBox_Information(nil, "Error", "No commands specified", qt.QMessageBox__Ok, qt.QMessageBox__Ok)
				return
			}

			resCh := make(chan string)
			errCh := make(chan string) // Stores the error message, in case of failed request
			go currentEnv.AddActuator(code.Text(), name.Text(), kind.Text(), commands, resCh, errCh)
			select {
			case res := <-resCh:
				qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			case error := <-errCh:
				qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			}
		})
		layout.AddWidget(addActBtn, 0, 0)
	}

	addSenBtn := qt.NewQPushButton2("Add Sensor", nil)
	addSenBtn.ConnectClicked(func(checked bool) {

		resCh := make(chan string)
		errCh := make(chan string) // Stores the error message, in case of failed request
		go currentEnv.AddSensor(code.Text(), name.Text(), kind.Text(), resCh, errCh)
		select {
		case res := <-resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(addSenBtn, 0, 0)


	doneBtn := qt.NewQPushButton2("Done", nil)
	doneBtn.ConnectClicked(func(checked bool) {

		resCh := make(chan string)
		errCh := make(chan string)
		go currentEnv.Done("add", resCh, errCh)
		select {
		case res := <-resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getHomepageWidget())
		case err := <-errCh:
			qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getHomepageWidget())
		}
	})
	layout.AddWidget(doneBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getSendCommandWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	var code string

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	envListCB := qt.NewQComboBox(nil)
	envListCB.SetEditable(false)
	selEnvBtn := qt.NewQPushButton2("Select Environment", nil)

	getEnvList(envListCB)

	// Gets all devices from server and shows them in a combo box
	devicesCB := qt.NewQComboBox(nil)
	devicesCB.SetEditable(false)
	devicesCB.SetVisible(false)

	selActBtn := qt.NewQPushButton2("Select Actuator", nil)
	selActBtn.SetVisible(false)

	selEnvBtn.ConnectClicked(func(checked bool) {

		client.SetCurrentEnv(currentEnv, envListCB.CurrentText())
		resCh := make(chan string)
		errCh := make(chan string)
		go currentEnv.GetActuatorList(resCh, errCh)

		select {
		case <- resCh:
			tmp := make([]string, 0)
			for _, value := range currentEnv.ActuatorMap {
				tmp = append(tmp, "Code: " + strconv.Itoa(value.Code) + " Name: " + value.Name)
			}
			fmt.Printf("TMP %v", tmp)
			devicesCB.Clear()
			devicesCB.AddItems(tmp)
			devicesCB.SetVisible(true)
			selActBtn.SetVisible(true)

		case err := <-errCh:
			qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}

	})

	commandsCB := qt.NewQComboBox(nil)
	commandsCB.SetEditable(false)
	commandsCB.SetVisible(false)
	sendBtn := qt.NewQPushButton2("Send", nil)
	sendBtn.SetVisible(false)

	selActBtn.ConnectClicked(func(checked bool) {

		// Extracts a number from the string
		re := regexp.MustCompile("[0-9]+")
		code = fmt.Sprintf(re.FindAllString(devicesCB.CurrentText(), -1)[0])

		commands := currentEnv.ActuatorMap[code].Commands
		commandsCB.Clear()
		commandsCB.AddItems(commands)
		commandsCB.SetVisible(true)
		sendBtn.SetVisible(true)

	})
	sendBtn.ConnectClicked(func(checked bool) {

		resCh := make(chan string)
		errCh := make(chan string)
		go currentEnv.SendCommand(code, commandsCB.CurrentText(), resCh, errCh)

		select {
		case res := <- resCh :
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		case err := <- errCh:
			qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})

	layout.AddWidget(envListCB, 0, 0)
	layout.AddWidget(selEnvBtn, 0, 0)
	layout.AddWidget(devicesCB, 0, 0)
	layout.AddWidget(selActBtn, 0, 0)
	layout.AddWidget(commandsCB, 0, 0)
	layout.AddWidget(sendBtn, 0, 0)
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getDeviceList(devicesCB * qt.QComboBox) {
	resCh := make(chan string)
	errCh := make(chan string)

	// Gets all environments from server and shows them in a combo box
	go currentEnv.GetDevicesList(resCh, errCh)

	select {
	case <- resCh:
		tmp := make([]string, 0)
		for _, value := range currentEnv.SensorMap {
			tmp = append(tmp, "Code: " + strconv.Itoa(value.Code) + " Name: " + value.Name + "[Sensor]")
		}
		for _, value := range currentEnv.ActuatorMap {
			tmp = append(tmp, "Code: " + strconv.Itoa(value.Code) + " Name: " + value.Name + "[Actuator]")
		}
		devicesCB.AddItems(tmp)

	case error := <-errCh:
		qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
	}
}

func getDeleteSensorsWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	devicesCB := qt.NewQComboBox(nil)
	devicesCB.SetEditable(false)

	getDeviceList(devicesCB)
	layout.AddWidget(devicesCB, 0, 0)

	code := qt.NewQLineEdit(nil)
	code.SetPlaceholderText("Code")
	layout.AddWidget(code, 0, 0)

	addSenBtn := qt.NewQPushButton2("Delete Sensor", nil)
	addSenBtn.ConnectClicked(func(checked bool) {
		resCh := make(chan string)
		errCh := make(chan string) // Stores the error message, in case of failed request
		go currentEnv.DeleteSensor(code.Text(), resCh, errCh)
		select {
		case res := <-resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(addSenBtn, 0, 0)

	addActBtn := qt.NewQPushButton2("Delete Actuator", nil)
	addActBtn.ConnectClicked(func(checked bool) {
		resCh := make(chan string)
		errCh := make(chan string) // Stores the error message, in case of failed request
		go currentEnv.DeleteActuator(code.Text(), resCh, errCh)
		select {
		case res := <-resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(addActBtn, 0, 0)

	doneBtn := qt.NewQPushButton2("Done", nil)
	doneBtn.ConnectClicked(func(checked bool) {
		resCh := make(chan string)
		errCh := make(chan string)
		go currentEnv.Done("delete", resCh, errCh)
		select {
		case res := <-resCh:
			qt.QMessageBox_Information(nil, "OK", res, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
			changeWindow(widget, getHomepageWidget())
		case err := <-errCh:
			qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(doneBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getEnvList(envListCB *qt.QComboBox) {
	envNameCh := make(chan string)
	errCh := make(chan string)

	// Gets all environments from server and shows them in a combo box
	go clientInstance.GetEnvList(envNameCh, errCh)

	select {
	case <- envNameCh:
		envList := make([] string, 0)
		for env := range envNameCh {
			envList = append(envList, env)
		}
		envListCB.AddItems(envList)
	case err := <-errCh:
		qt.QMessageBox_Information(nil, "Error", err, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
	}
}

func getHistoryDataWidget() *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	envListCB := qt.NewQComboBox(nil)
	envListCB.SetEditable(false)

	getEnvList(envListCB)
	layout.AddWidget(envListCB, 0, 0)

	selEnvBtn := qt.NewQPushButton2("Select Environment", nil)
	selEnvBtn.ConnectClicked(func(checked bool) {
		client.SetCurrentEnv(currentEnv, envListCB.CurrentText())
		resCh := make(chan [] * client.Measurement)
		errCh := make(chan string) // Stores the error message, in case of failed request
		go currentEnv.GetHistoryData(resCh, errCh)
		select {
		case history := <- resCh:
			changeWindow(widget, getVisualizeHistoryWidget(history))
		case error := <-errCh:
			qt.QMessageBox_Information(nil, "Error", error, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
		}
	})
	layout.AddWidget(selEnvBtn, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}

func getVisualizeHistoryWidget(history [] * client.Measurement) *qt.QWidget {

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	dataList := qt.NewQListWidget(nil)
	dataList.AddItem(fmt.Sprintf("Code\tValue\tTimestamp\n"))

	for _, sensorDataRow := range history {
		dataList.AddItem(fmt.Sprintf("%d\t%f\t%s\n", sensorDataRow.Code, sensorDataRow.Value, sensorDataRow.Timestamp))
	}

	layout.AddWidget(dataList, 0, 0)

	backBtn := qt.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(checked bool) {
		changeWindow(widget, getHomepageWidget())
	})
	layout.AddWidget(backBtn, 0, 0)

	return widget

}


func main() {

	currentEnv = client.NewEnvironment()
	clientInstance = client.GetClientInstance()
	initWindow()
	qt.QApplication_Exec()

}
