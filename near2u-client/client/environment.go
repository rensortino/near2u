package client

import (
	"../utils"
	"strconv"
)

type Environment struct {
	Name      string                 `json:"name"`
	SensorMap map[string]interface{} `json:"sensors"`
	LastModified int `json:"lastmodified"`
}

type Device struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Actuator struct {
	Device `json:"device"`
	Commands []string `json:"commands"`
}

// Temporary list that stores the sensors to add / delete
var sensorList [] Sensor

func (e * Environment) GetSensorList(sensorListCh chan []string, errCh chan string) {

	data := struct {
		EnvName string
	} {
		e.Name,
	}

	res := utils.SocketCommunicate("visualizza_sensori", clientInstance.LoggedUser, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	sensors := res["data"].(map[string]interface{})["sensors"].([]Sensor)

	for _, sensor := range sensors {
		e.SensorMap[strconv.Itoa(sensor.Code)] = sensor
	}

	// TODO Make the server return a list of sensors
	// Returns a list of sensor codes
	sensorListCh <- res["data"].(map[string]interface{})["sensors"].([]string)
}

func (e * Environment) AddSensor(code, name, kind string,
	sensorCh chan interface{}, errCh chan string) {

	if !(len(e.SensorMap) == 0) {
		_, found := e.SensorMap[code]
		if found {
			errCh <- "Code already in use, please choose another one"
			close(errCh)
			return
		}
	}
	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Sensor Code must be integer"
		close(errCh)
		return
	}
	newSensor := Sensor{
		intCode,
		name,
		kind,
	}
	sensorList = append(sensorList, newSensor)
	sensorCh <- &newSensor
	close(sensorCh)
}

// Adds the selected sensor to a list in order to collect all sensors to delete
// and to send them in batch to the server
func (e * Environment) DeleteSensor(code string, sensorCh chan interface{}, errCh chan string) {

	if len(e.SensorMap) == 0 {
		errCh <- "Empty Environment, no sensors to delete"
		close(errCh)
		return
	}
	sensor, found := e.SensorMap[code]
	if found {
		// Sensor list is passed by reference because is a slice
		sensorList = append(sensorList, sensor.(Sensor))
		sensorCh <- sensor.(Sensor)
		close(sensorCh)
		return
	} else {
		errCh <- "Sensor not found"
		close(errCh)
		return
	}

}

func (e * Environment) Done(resCh, errCh chan string) {

	if len(sensorList) == 0 {
		errCh <- "No sensor selected"
		return
	}

	data := struct {
		Sensors []Sensor `json:"sensors"`
		EnvName string   `json:"envName"`
	}{
		sensorList,
		e.Name,
	}

	res := utils.SocketCommunicate("inserisci_sensori", clientInstance.LoggedUser, data)

	if res["status"] == "Succesfull" {
		resCh <- res["status"].(string)
		close(resCh)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
}