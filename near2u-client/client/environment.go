package client

import (
	"../utils"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
)

type Environment struct {
	Name         string
	SensorMap    map[string]Sensor
	ActuatorMap    map[string]Actuator
	LastModified int
}

// Temporary list that stores the devices to add / delete
// TODO Check utility of deviceList as interface{}, consider to convert it to []int
var deviceList []interface{} // used list of interfaces to implement polymorphism

var deviceToDelete []int

func (e * Environment) GetSensorsList(resCh chan string, errCh chan string) {

	data := struct {
		EnvName string `json:"envname"`
		Type string `json:"type"`
	}{
		e.Name,
		"sensors",
	}

	res := utils.SocketCommunicate("visualizza_dispositivi", clientInstance.LoggedUser.Auth, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	if len(res["data"].(map[string]interface{})["devices"].([]interface{})) == 0 {
		resCh <- "Environment with no Sensors"
		return
	}

	sensors := res["data"].(map[string]interface{})["devices"].([]interface{})

	for _, sensor := range sensors {

		s := sensor.(map[string]interface{})

		tmpSensor := &Sensor{
			Device{
				int(s["code"].(float64)),
				s["name"].(string),
				s["kind"].(string),
			},
			0.0,
		}
		e.SensorMap[strconv.Itoa(tmpSensor.Code)] = *tmpSensor
	}

	resCh <- "Success"
	//close(resCh)
}

func (e * Environment) GetActuatorList(resCh chan string, errCh chan string) {

	data := struct {
		EnvName string `json:"envname"`
		Type string `json:"type"`
	}{
		e.Name,
		"actuators",
	}

	res := utils.SocketCommunicate("visualizza_dispositivi", clientInstance.LoggedUser.Auth, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	if len(res["data"].(map[string]interface{})["devices"].([]interface{})) == 0 {
		resCh <- "Environment with no Actuators"
		return
	}

	actuators := res["data"].(map[string]interface{})["devices"].([]interface{})

	for _, act := range actuators {

		a := act.(map[string]interface{})

		commands := make([]string, 0)

		for _, cmd := range a["commands"].([]interface{}) {
			commands = append(commands, cmd.(string))
		}

		tmpAct := &Actuator{
			Device{
				int(a["code"].(float64)),
				a["name"].(string),
				a["kind"].(string),
			},
			commands,
		}
		e.ActuatorMap[strconv.Itoa(tmpAct.Code)] = *tmpAct
	}

	resCh <- "Success"
	close(resCh)
}

func (e * Environment) GetDevicesList(resCh chan string, errCh chan string) {

	e.GetSensorsList(resCh, errCh)

	e.GetActuatorList(resCh, errCh)
}

func (e * Environment) AddSensor(code, name, kind string,
	resCh, errCh chan string) {

	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Code must be integer"
		close(errCh)
		return
	}

	if !(len(e.SensorMap) == 0){
		_, found := e.SensorMap[string(intCode)]
		if found {
			errCh <- "Code already in use, please choose another one"
			close(errCh)
			return
		}
	}
	var newDevice = NewDevice(intCode, name, kind)
	var ok bool // stores the value of the append function

	newSensor := NewSensor(newDevice)
	deviceList, ok = newSensor.Append(deviceList)
	if ok { // Append returns true if the operation succeeded
		resCh <- fmt.Sprintf("Sensor %s added to the list", newSensor.Name)
		close(resCh)
	}
}

func (e * Environment) AddActuator(code, name, kind string, commands []string,
	resCh, errCh chan string) {

	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Code must be integer"
		close(errCh)
		return
	}

	if !(len(e.ActuatorMap) == 0){
		_, found := e.ActuatorMap[string(intCode)]
		if found {
			errCh <- "Code already in use, please choose another one"
			close(errCh)
			return
		}
	}
	var newDevice = NewDevice(intCode, name, kind)
	var ok bool // stores the value of the append function

	newActuator := NewActuator(newDevice, commands)
	deviceList, ok = newActuator.Append(deviceList)
	if ok { // Append returns true if the operation succeeded
		resCh <- fmt.Sprintf("Actuator %s added to the list", newActuator.Name)
		close(resCh)
	}
}

// Adds the selected sensor to a list in order to collect all sensors to delete
// and to send them in batch to the server
func (e * Environment) DeleteSensor(code string, resCh, errCh chan string) {

	if len(e.SensorMap) == 0 {
		errCh <- "No sensor to delete"
		close(errCh)
		return
	}
	s, found := e.SensorMap[code]
	if found {
		deviceToDelete = append(deviceToDelete, s.Code)
		resCh <- "Sensor added for deletion"
	} else {
		errCh <- "Sensor not found"
		close(errCh)
		return
	}

}

func (e * Environment) DeleteActuator(code string, resCh, errCh chan string) {

	if len(e.ActuatorMap) == 0 {
		errCh <- "No actuator to delete"
		close(errCh)
		return
	}
	a, found := e.ActuatorMap[code]
	if found {
		deviceToDelete = append(deviceToDelete, a.Code)
		resCh <- "Actuator added for deletion"

	} else {
		errCh <- "Actuator not found"
		close(errCh)
		return
	}

}

func (e * Environment) Done(operation string, resCh, errCh chan string) {


	var res map[string]interface{}

	switch operation {
	case "add":
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
		res = utils.SocketCommunicate("inserisci_dispositivi", clientInstance.LoggedUser.Auth, data)
		deviceList = make([]interface{}, 0) // Empties the list for future requests
	case "delete":
		if len(deviceToDelete) == 0 {
			errCh <- "No device selected"
			return
		}

		data := struct {
			Devices []int `json:"devices"`
			EnvName string   `json:"envname"`
		}{
			deviceToDelete,
			e.Name,
		}
		res = utils.SocketCommunicate("elimina_dispositivi", clientInstance.LoggedUser.Auth, data)
		deviceToDelete = make([]int, 0) // Empties the list for future requests
	}


	if res["status"] == "Successful" {
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

	act, found := e.ActuatorMap[code]

	if !found {
		errCh <- "Actuator not found"
		close(errCh)
		return
	}

	if act.Commands == nil {
		errCh <- "Actuator has no commands"
	}

	res := act.SendCommand(e.Name, command)

	if res != "error" {
		resCh <- res
		return
	} else {
		errCh <- "Error sending commands"
	}
}

func (e * Environment) GetHistoryData(resCh chan [] * Measurement, errCh chan string)  {

	data := struct {
		EnvName string `json:"envname"`
	} {
		e.Name,
	}

	res := utils.SocketCommunicate("visualizza_storico", clientInstance.LoggedUser.Auth, data)


	if res["status"] == "Successful" {
		historyFile, err := os.Create("history.csv")
		if err != nil {
			log.Fatalln(err)
		}
		defer historyFile.Close()
		historyFile.WriteString(fmt.Sprintf("Code,Value,Timestamp\n"))
		sensorData := res["data"].(map[string]interface{})["sensor_data"].([]interface{})
		history := make([] * Measurement, 0)

		for _, data := range sensorData {

			d := data.(map[string]interface{})

			measurement := &Measurement{
				int(d["code"].(float64)),
				d["misura"].(float64),
				d["time"].(string),
			}
			history = append(history, measurement)
			historyFile.WriteString(fmt.Sprintf("%d,%f,%s\n", measurement.Code, measurement.Value, measurement.Timestamp))
		}

		resCh <- history
		close(resCh)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
}

func (e * Environment) GetPlot() image.Image {

	plotFile, err := os.Open("grafico.png")
	if err != nil {
		log.Fatalln("Error opening plot file")
	}
	defer plotFile.Close()

	img, err := png.Decode(plotFile)
	if err != nil {
		log.Fatalln("Error decoding image")
	}

	return img
}