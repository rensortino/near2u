package client

import (
	"../utils"
	"strconv"
)

type Environment struct {
	Name         string                 `json:"name"`
	DeviceMap    map[string]interface{} `json:"devices"`
	LastModified int                    `json:"lastmodified"`
}

// Temporary list that stores the devices to add / delete
var deviceList []interface{} // used list of interfaces to implement polymorphism

func (e * Environment) GetDevicesList(deviceListCh chan []interface{}, errCh chan string) {

	data := struct {
		EnvName string
	}{
		e.Name,
	}

	res := utils.SocketCommunicate("visualizza_dispositivi", clientInstance.LoggedUser, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	devices := res["data"].(map[string]interface{})["sensors"].([]interface{})

	for _, device := range devices {
		e.DeviceMap[strconv.Itoa(device.(Device).Code)] = device
	}

	// Returns a list of devices
	deviceListCh <- devices
}

func (e * Environment) AddDevice(code, name, kind string, commands []string,
	sensorCh chan interface{}, errCh chan string) {

	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Sensor Code must be integer"
		close(errCh)
		return
	}

	newDevice := NewDevice(intCode, name, kind, commands)

	if newDevice == nil {
		errCh <- "Error creating new device"
		close(errCh)
		return
	}

	var mapCode int

	if !(len(e.DeviceMap) == 0) {
		switch d := newDevice.(type){
		case Sensor:
			mapCode = d.Code
		case Actuator:
			mapCode = d.Code
		}
		_, found := e.DeviceMap[string(mapCode)]
		if found {
			errCh <- "Code already in use, please choose another one"
			close(errCh)
			return
		}
	}

	switch d := newDevice.(type){
	case Sensor:
		d.Append(deviceList)
	case Actuator:
		d.Append(deviceList)
	}

	sensorCh <- &newDevice
	close(sensorCh)
}

// Adds the selected sensor to a list in order to collect all sensors to delete
// and to send them in batch to the server
func (e * Environment) DeleteDevice(code string, deviceCh chan interface{}, errCh chan string) {

	if len(e.DeviceMap) == 0 {
		errCh <- "Empty Environment, no sensors to delete"
		close(errCh)
		return
	}
	device, found := e.DeviceMap[code]
	if found {
		// Sensor list is passed by reference because is a slice
		switch d := device.(type) {
		case Sensor:
			d.Append(deviceList)
		case Actuator:
			d.Append(deviceList)
		}
		deviceCh <- device
		close(deviceCh)
		return
	} else {
		errCh <- "Sensor not found"
		close(errCh)
		return
	}

}

func (e * Environment) Done(operation string, resCh, errCh chan string) {

	if len(deviceList) == 0 {
		errCh <- "No device selected"
		return
	}

	data := struct {
		Devices []interface{} `json:"devices"`
		EnvName string   `json:"envname"`
	}{
		deviceList,
		e.Name,
	}

	var res map[string]interface{}

	switch operation {
	case "add":
		res = utils.SocketCommunicate("inserisci_dispositivi", clientInstance.LoggedUser, data)
	case "delete":
		res = utils.SocketCommunicate("elimina_dispositivi", clientInstance.LoggedUser, data)
	}


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

func (e * Environment) SendCommand(code, command string, resCh, errCh chan string) {

	act, found := e.DeviceMap[code]

	if !found {
		errCh <- "Actuator not found"
		close(errCh)
		return
	}

	if act.(Actuator).Commands == nil {
		errCh <- "Actuator has no commands"
	}

	res := e.DeviceMap[code].(Actuator).SendCommand(e.Name, command)

	if res != "error" {
		resCh <- res
		close(resCh)
		return
	} else {
		errCh <- "Error sending commands"
		close(errCh)
	}
}