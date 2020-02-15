package client

import (
	"../utils"
	"encoding/json"
	"net"
	"strconv"
)

type Sensor struct {
	ID string `json:ID`
	Name string `json:name`
	Measurement float32 `json:measurement`
}

type Environment struct {
	ID string `json:ID`
	SensorMap map[string]Sensor `json:sensors`
}

type SelEnvRequest struct {
	EnvID string `json:"env"`
	utils.RequestParams `json:"params"`
}

// TODO Create Client struct


// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func GetEnvList(conn net.Conn) []string {

	/*
	// TODO use client token
	request := utils.RequestParams {
		"getEnvList",
		"auth",
	}

	jsonReq, _ := json.Marshal(request)

	utils.SocketSend(conn, jsonReq)

	rx := make(chan []byte)
	utils.SocketReceive(conn, rx)

	// TODO Take environment list accessing the json
	var envList struct {
		Environments [] string `json:"environments"`
	}
	json.Unmarshal(<- rx, &envList)
	log.Println(envList.Environments)

	 */

	// TODO Delete test string
	var test = make([]string, 10)
	for i := 0; i < 10; i++ {
		test[i] = "Test " + strconv.Itoa(i)
	}
	return test
}

func getSensorData() {
	// TODO Subscribe to MQTT
}

func SelectEnv(conn net.Conn, rx chan []byte, envID string) {

	params := utils.RequestParams {
		"getEnvList",
		"auth",
	}

	// TODO Make normal struct
	request := struct {
		utils.RequestParams
		string `json:"envID"`
	}{
		params,
		envID,
	}

	jsonReq, _ := json.Marshal(request)

	utils.SocketSend(conn, jsonReq)
	//rx := make(chan [] byte)
	// Returns sensor list
	//utils.SocketReceive(conn, rx)

	// receivedData := <- rx

	// TODO Define sensor map, {sensorID, measurement}
	// Map needed to update values coming from MQTT broker directly using ID
	// e.g. sensorMap[ID] = sensorData
	// sensorMap = Json.Decode(receivedData)["data"]
	// return sensorList

}