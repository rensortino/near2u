package client

import "../utils"

type Device struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type DeviceI interface {
	Append(deviceList []interface{}) interface{}
}

type Actuator struct {
	Device `json:"device"`
	Commands []string `json:"commands"`
}

type Sensor struct {
	Device `json:"device"`
	Measurement float64 `json:"measurement"`
}

type Measurement struct {
	Code int `json:"code"`
	Value float64 `json:"misura"`
	Timestamp string `json:"time"`
}

func NewDevice(code int, name, kind string) Device {

	newDevice := Device {
		code,
		name,
		kind,
	}
	return newDevice
}

func NewSensor(dev Device) * Sensor {
	newSensor := &Sensor{
		dev,
		0.0,
	}
	return newSensor
}

func NewActuator(dev Device, commands []string) * Actuator {

	newActuator := &Actuator{
		dev,
		commands,
	}
	return newActuator
}

func (s * Sensor) Append(deviceList []interface{}) ([]interface{}, bool) {

	// The following struct is needed to flatten the JSON before transmitting it to the server
	sensor := struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		Kind string `json:"kind"`
	} {
		s.Code,
		s.Name,
		s.Kind,
	}

	for _, dev := range deviceList {
		switch d := dev.(type) {
		case Sensor:
			if d.Code == s.Code {
				return deviceList, false
			}
		case Actuator:
			if d.Code == s.Code {
				return deviceList, false
			}
		}
	}
	return append(deviceList, sensor), true
}

func (a * Actuator) Append(deviceList []interface{}) ([]interface{}, bool) {

	// The following struct is needed to flatten the JSON before transmitting it to the server
	actuator := struct {
		Code int    `json:"code"`
		Name string `json:"name"`
		Kind string `json:"kind"`
		Commands []string `json:"commands"`
	} {
		a.Code,
		a.Name,
		a.Kind,
		a.Commands,
	}

	for _, dev := range deviceList {
		switch d := dev.(type) {
		case Sensor:
			if d.Code == a.Code {
				return deviceList, false
			}
		case Actuator:
			if d.Code == a.Code {
				return deviceList, false
			}
		}
	}

	return append(deviceList, actuator), true
}

func (a * Actuator) SendCommand(envName, command string) string {
	for _, comm := range a.Commands {
		if command == comm {
			data := struct {
				EnvName string `json:"envname"`
				ActCode int `json:"code"`
				Command string `json:"command"`
			} {
				envName,
				a.Code,
				command,
			}
			res := utils.SocketCommunicate("invia_comando", clientInstance.LoggedUser.Auth, data)

			if res["status"] == "Successful" {
				return res["data"].(string)
			}
		}
	}
	return "error"
}