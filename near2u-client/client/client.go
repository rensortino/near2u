package client

import (
	"../utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net"
	"strconv"
)


type Sensor struct {
	ID string
	Name string
	Measurement float32
}

type Environment struct {
	ID string
	sensorMap map[string]Sensor
}

var (
	function string
	data string
	auth string
)

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func GetEnvList(conn net.Conn) []string {
/*
	function = "getEnvList"
	data = ""
	auth = "a"

	socketSend(conn, function, data, auth)
	SocketReceive(conn)
*/
	//return strings.Split(<- Socket, ";")
	// TODO Delete test string
	var test = make([]string, 10)
	for i := 0; i < 10; i++ {
		test[i] = "Test " + strconv.Itoa(i)
	}
	return test
}

func getSensorData() {
	// Subscribe to MQTT
}

func selectEnv(conn net.Conn, rx chan string, envID string) {

	function = "selectEnvironment"
	data = envID

	utils.SocketSend(conn, function, data, auth)
	// Returns sensor list
	utils.SocketReceive(conn, rx)

	// receivedData := <- rx

	// TODO Define sensor map, {sensorID, measurement}
	// Map needed to update values coming from MQTT broker directly using ID
	// e.g. sensorMap[ID] = sensorData
	// sensorMap = Json.Decode(receivedData)["data"]
	// return sensorList

}