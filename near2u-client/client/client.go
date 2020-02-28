package client

import (
	"encoding/json"
	"log"
	"strconv"

	"../utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	ID         string
	LoggedUser string
	MQTTClient mqtt.Client
}

type Sensor struct {
	Code int    `json:"code"`
	Name string `json:name`
	Kind string `json:kind`
}

type Environment struct {
	Name      string                 `json:name`
	SensorMap map[string]interface{} `json:sensors`
}

var clientInstance *Client

// Implements singleton pattern
func GetClientInstance() *Client {

	if clientInstance == nil {
		clientInstance = &Client{
			"ID1",
			"",
			nil,
		}
	}

	return clientInstance
}

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func (c *Client) GetEnvList(envListCh chan []string, errCh chan string) {

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate("visualizza_ambienti", c.LoggedUser, nil, rx)

	res := <-rx
	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	envListCh <- res["data"].(map[string]interface{})["environments"].([]string)
	close(envListCh)
}

func (c *Client) GetSensorList(envName string, sensorListCh chan []string, errCh chan string) {

	data := struct {
		EnvName string
	} {
		envName,
	}

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate("visualizza_sensori", c.LoggedUser, data, rx)

	res := <-rx
	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	// Returns a list of sensor codes
	sensorListCh <- res["data"].(map[string]interface{})["sensors"].([]string)
}

func (c *Client) GetSensorData(topic string, rtCh chan map[string]interface{}) {

	c.MQTTClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Executes every time a message is published on the topic
		sensors := make(map[string]interface{})
		//env := Environment{SensorMap:sensors}

		json.Unmarshal(msg.Payload(), &sensors)

		log.Println("Data Received")
		log.Println(sensors)
		/*
			 Payload format:
			{"ID":"env1","SensorMap":{
				"sensor1":{"ID":"id","Name":"name","Measurement":7.4 },
				"sensor2":{"ID":"otherID","Name":"name2","Measurement":4.76}
			}}
		*/

		rtCh <- sensors["SensorMap"].(map[string]interface{})
	})
}

// Gracefully stops getting data from the broker
func (c *Client) StopGettingData(topic string, rtCh chan map[string]interface{}, quit chan bool) {
	c.MQTTClient.Unsubscribe(topic)
	// Empties the channel before closing it
	select {
	case <-rtCh:
		close(rtCh)
	default:
		close(rtCh)
	}
	quit <- true
	close(quit)
}

func (c *Client) GetTopicAndUri(envName string, topicCh, uriCh, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	rx := make(chan map[string]interface{})

	//Returns broker's address on rx channel
	go utils.SocketCommunicate("topic_ambiente", c.LoggedUser, data, rx)

	res := <-rx
	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	topicCh <- res["data"].(map[string]interface{})["topic"].(string)
	uriCh <- res["data"].(map[string]interface{})["broker_host"].(string)
	close(topicCh)
	close(uriCh)
}

func (c * Client) SelectEnv(envName string, envCh chan * Environment, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	rx := make(chan map[string]interface{})

	//Returns broker's address on rx channel
	go utils.SocketCommunicate("seleziona_ambiente", c.LoggedUser, data, rx)

	res := <-rx
	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	envCh <- res["data"].(map[string]interface{})["environment"].(* Environment)
	close(envCh)
}

func (c *Client) CreateEnv(envName string, envCh chan *Environment, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate("configura_ambiente", c.LoggedUser, data, rx)

	res := <-rx

	if res["status"] == "Succesfull" {
		newEnv := &Environment{
			envName,
			make(map[string]interface{}),
		}
		envCh <- newEnv
	} else {
		errCh <- res["error"].(string)
	}
}

func (c *Client) AddSensor(code, name, kind string, env *Environment,
	sensorCh chan interface{}, errCh chan string) {

	if !(len(env.SensorMap) == 0) {
		_, found := env.SensorMap[code]
		if found {
			errCh <- "Code already in use, please choose another one"
			return
		}
	}
	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Sensor Code must be integer"
		return
	}
	newSensor := Sensor{
		intCode,
		name,
		kind,
	}
	env.SensorMap[code] = newSensor
	sensorCh <- &newSensor

}

func (c *Client) DeleteSensor(code string, env *Environment, sensorList []Sensor, sensorCh chan interface{}, errCh chan string) {

	if len(env.SensorMap) == 0 {
		errCh <- "Empty Environment, no sensors to delete"
		return
	}
	sensor, found := env.SensorMap[code]
	if found {
		sensorList = append(sensorList, sensor.(Sensor))
		sensorCh <- sensor.(Sensor)
	} else {
		errCh <- "Sensor not found"
	}

}

func (c *Client) Done(envName string, sensorList [] Sensor, resCh, errCh chan string) {

	data := struct {
		Sensors []Sensor `json:"sensors"`
		EnvName string   `json:"envName"`
	}{
		sensorList,
		envName,
	}

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate("inserisci_sensori", c.LoggedUser, data, rx)

	res := <-rx

	if res["status"] == "Succesfull" {
		resCh <- res["status"].(string)
	} else {
		errCh <- res["error"].(string)
	}
}
