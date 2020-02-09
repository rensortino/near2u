package gui

import (
	"../client"
	qtcore "github.com/therecipe/qt/core"
	qtgui "github.com/therecipe/qt/gui"
	qt "github.com/therecipe/qt/widgets"
	"os"
)

var (
	mainWindow * qt.QMainWindow
	mainWidget * qt.QWidget
)

func getHomepageWidget() * qt.QWidget{

	layout := qt.NewQVBoxLayout()

	widget := qt.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	loginBtn := qt.NewQPushButton2("Log In", nil)
	loginBtn.ConnectClicked(func (checked bool) {
		widget.Destroy(true,  true)
		mainWindow.SetCentralWidget(getLoginWidget())
	})
	layout.AddWidget(loginBtn, 0, 0)

	registerBtn := qt.NewQPushButton2("Register", nil)
	registerBtn.ConnectClicked(func (checked bool) {
		widget.Destroy(true,  true)
		mainWindow.SetCentralWidget(getRegisterWidget())
	})
	layout.AddWidget(registerBtn, 0, 0)

	selEnvBtn := qt.NewQPushButton2("Select Environment", nil)
	selEnvBtn.ConnectClicked(func (checked bool) {
	})
	layout.AddWidget(selEnvBtn, 0, 0)

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

	rx := qtcore.NewQRegExp()
	rx.SetPattern("\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,4}\\b")
	layout.AddWidget(email, 0, 0)
	validator := qtgui.NewQRegExpValidator2(rx, nil)
	email.SetValidator(validator)

	password := qt.NewQLineEdit(nil)
	password.SetPlaceholderText("Password")
	layout.AddWidget(password, 0, 0)

	registerBtn := qt.NewQPushButton2("Register", nil)
	registerBtn.ConnectClicked(func(checked bool) {
		if conn != nil {
			client.Register(conn, name.Text(), surname.Text(), email.Text(), password.Text())
		}
	})
	layout.AddWidget(registerBtn, 0, 0)

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

	button := qt.NewQPushButton2("Log In", nil)
	button.ConnectClicked(func(checked bool) {
		if conn != nil {
			client.Login(conn, email.Text(), password.Text())
		}
	})
	layout.AddWidget(button, 0, 0)

	return widget

}

func InitWindow() {

	// Creates a new graphic application
	qt.NewQApplication(len(os.Args), os.Args)

	mainWindow = qt.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle("Near2U")
	mainWindow.SetMinimumSize2(400, 400)

	mainWidget = getHomepageWidget()

	mainWindow.SetCentralWidget(mainWidget)
	mainWindow.Show()
}