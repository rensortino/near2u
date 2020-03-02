package client

import "../utils"

type Device struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type DeviceI interface {
	Append(deviceList []interface{}) interface{}
	Delete(code string) interface{}
}

type Actuator struct {
	Device `json:"device"`
	Commands []string `json:"commands"`
}

type Sensor struct {
	Device `json:"device"`
	Measurement float64 `json:"measurement"`
}

func NewDevice(code int, name, kind string, commands []string) interface{} {

	newDevice := Device {
		code,
		name,
		kind,
	}
	if commands == nil {
		newSensor := Sensor{
			newDevice,
			nil,
		}
		return newSensor
	}
	if commands != nil {
		newActuator := Actuator {
			newDevice,
			commands,
		}
		return newActuator
	}

	return nil
}

func (s * Sensor) Append(deviceList []interface{}) {

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

	deviceList = append(deviceList, sensor)
}

func (a * Actuator) Append(deviceList []interface{}) {

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

	deviceList = append(deviceList, actuator)
}

func (a * Actuator) SendCommand(envName, command string) string {
	for _, comm := range a.Commands {
		if command == comm {
			data := struct {
				EnvName string `json:"envname"`
				Command string `json:"command"`
			} {
				envName,
				command,
			}
			res := utils.SocketCommunicate("invia_comando", clientInstance.LoggedUser, data)

			if res["status"] == "Succesfull" {
				return res["result"].(string)
			}
		}
	}
	return "error"
}