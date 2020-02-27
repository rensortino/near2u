package client

import (
	"../utils"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"strconv"
)

type Client struct {
	ID string
	LoggedUser string
	MQTTClient mqtt.Client
}

type Sensor struct {
	Code int `json:"code"`
	Name string `json:name`
	Kind string `json:kind`
}

type Environment struct {
	Name string `json:name`
	SensorMap map[string]interface{} `json:sensors`
}

var clientInstance * Client

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
func (c * Client) GetEnvList() []string {

	/*
	request := struct {
		Function string
		Auth string
	} {
		"getEnvList"
		clientInstance.LoggedUser
	}

	jsonReq, _ := json.Marshal(request)

	rx := make(chan []byte)
	go SocketCommunicate(jsonReq, rx)

	var envList struct {
	Environments [] string `json:"environments"`
	} {}

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

func (c * Client) GetSensorData(topic string, rtCh chan map[string]interface{}) {

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
func (c * Client) StopGettingData(topic string, rtCh chan map[string]interface{}, quit chan bool) {
	c.MQTTClient.Unsubscribe(topic)
	// Empties the channel before closing it
	select {
		case <- rtCh:
			close(rtCh)
		default:
			close(rtCh)
	}
	quit <- true
	close(quit)
}

func (c * Client) SelectEnv(envName string, topicCh, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	rx := make(chan map[string]interface{})

	//Returns broker's address on rx channel
	go utils.SocketCommunicate("seleziona_ambiente", c.LoggedUser, data, rx)

	res := <- rx
	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	uri, err := url.Parse("tcp://" + res["data"].(map[string]interface{})["broker_host"].(string))
	if err != nil {
		log.Fatal(err)
	}

	if c.MQTTClient == nil {
		c.MQTTClient = utils.MQTTConnect(c.ID, uri)
	}

	log.Println(res["data"].(map[string]interface{})["topic"])

	topicCh <- res["data"].(map[string]interface{})["topic"].(string)
	close(topicCh)
}

func (c * Client) CreateEnv(envName string, envCh chan * Environment, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	} {
		envName,
	}

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate("configura_ambiente", c.LoggedUser, data, rx)

	res := <- rx

	if res["status"] ==  "Succesfull" {
		newEnv := &Environment {
			envName,
			make(map[string]interface{}),
		}
		envCh <- newEnv
	} else {
		errCh <- res["error"].(string)
	}
}

func (c * Client) AddSensor(code, name, kind string, newEnv * Environment,
							sensorCh chan interface{}, errCh chan string) {

	if !(len(newEnv.SensorMap) == 0) {
		log.Println(len(newEnv.SensorMap))
		log.Println("Printing newEnv.SensorMap[code]")
		log.Println(newEnv.SensorMap[code])
		_, found := newEnv.SensorMap[code]
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
	newSensor := Sensor {
		intCode,
		name,
		kind,
	}
	newEnv.SensorMap[code] = newSensor
	sensorCh <- &newSensor

}

func (c * Client) Done(newEnv * Environment, resCh, errCh chan string) {

	sensorList := make([]Sensor, 0)

	for _, sensor := range newEnv.SensorMap {
		sensorList = append(sensorList, sensor.(Sensor))
		log.Println(sensor)
	}

	data := struct {
		Sensors []Sensor `json:"sensors"`
		EnvName string `json:"envName"`
	} {
		sensorList,
		newEnv.Name,
	}

	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate("inserisci_sensori", c.LoggedUser, data, rx)

	res := <- rx

	if res["status"] ==  "Succesfull" {
		resCh <- res["status"].(string)
	} else {
		errCh <- res["error"].(string)
	}
}