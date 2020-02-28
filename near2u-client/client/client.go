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
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Environment struct {
	Name      string                 `json:"envname"`
	SensorMap map[int]interface{} `json:"sensors"`
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

	environments := make([]string, 0)

	for _, env := range res["data"].(map[string]interface{})["environments"].([]interface{}) {
		log.Println(env)
		environments = append(environments, env.(string))
	}
	envListCh <- environments
	close(envListCh)
}

func (c *Client) GetSensorList(envName string, sensorListCh chan []Sensor, errCh chan string) {

	data := struct {
		EnvName string `json:"envname"`
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
	sensorList := res["data"].(map[string]interface{})["sensors"]
	sensors := make([] Sensor, 0)

	for _, sensorJSON := range sensorList.([]interface{}) {
		sensor := &Sensor{}
		sensorData, _ := json.Marshal(sensorJSON)
		json.Unmarshal(sensorData, sensor)
		sensors = append(sensors, * sensor)
	}
	sensorListCh <- sensors
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

// TODO Change sensorlist type to []Sensor
func (c * Client) SelectEnv(envName string, envCh chan * Environment, errCh chan string) {

	sensorListCh := make(chan [] Sensor)

	go c.GetSensorList(envName, sensorListCh, errCh)

	env := &Environment{
		envName,
		make(map[int]interface{}),
	}

	select {
	case res := <-sensorListCh:
		// TODO Change with sensor code
		for _, sensor := range res {
			env.SensorMap[sensor.Code] = sensor
			log.Println(env.SensorMap[sensor.Code])
		}
	case err := <- errCh:
		log.Println(err)
		errCh <- err
		return
	}
	

	envCh <- env
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
			make(map[int]interface{}),
		}
		envCh <- newEnv
	} else {
		errCh <- res["error"].(string)
	}
}

// TODO Change sensorlist type to []Sensor
func (c *Client) AddSensor(code, name, kind string, env *Environment,
	sensorCh chan interface{}, errCh chan string) {

	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Sensor Code must be integer"
		return
	}

	if !(len(env.SensorMap) == 0) {
		_, found := env.SensorMap[intCode]
		if found {
			errCh <- "Code already in use, please choose another one"
			return
		}
	}
	
	newSensor := Sensor{
		intCode,
		name,
		kind,
	}
	env.SensorMap[intCode] = newSensor
	sensorCh <- &newSensor

}

// TODO Change sensorlist type to []Sensor
func (c *Client) DeleteSensor(code string, env *Environment, sensorList []Sensor, sensorCh chan interface{}, errCh chan string) {

	intCode, err := strconv.Atoi(code)
	if err != nil {
		errCh <- "Sensor Code must be integer"
		return
	}

	if len(env.SensorMap) == 0 {
		errCh <- "Empty Environment, no sensors to delete"
		return
	}
	sensor, found := env.SensorMap[intCode]
	if found {
		sensorList = append(sensorList, sensor.(Sensor))
		sensorCh <- sensor.(Sensor)
	} else {
		errCh <- "Sensor not found"
	}

}

func (c *Client) Done(envName, function string, sensorList [] Sensor, resCh, errCh chan string) {

	data := struct {
		Sensors []Sensor `json:"sensors"`
		EnvName string   `json:"envname"`
	}{
		sensorList,
		envName,
	}

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate(function, c.LoggedUser, data, rx)

	res := <-rx

	if res["status"] == "Succesfull" {
		resCh <- res["status"].(string)
	} else {
		errCh <- res["error"].(string)
	}
}
