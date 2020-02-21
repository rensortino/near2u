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
	ID string `json:ID`
	Name string `json:name`
	Measurement float32 `json:measurement`
}

type Environment struct {
	ID string `json:ID`
	SensorMap map[string]Sensor `json:sensors`
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

func (c * Client) GetSensorData(topic string, rtCh chan map[string]Sensor) {

	c.MQTTClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Executes every time a message is published on the topic
		sensors := make(map[string]Sensor)
		env := Environment{SensorMap:sensors}

		json.Unmarshal(msg.Payload(), &env)
		/*
			 Payload format:
			{"ID":"env1","SensorMap":{
				"sensor1":{"ID":"id","Name":"name","Measurement":7.4 },
				"sensor2":{"ID":"otherID","Name":"name2","Measurement":4.76}
			}}
		*/

		rtCh <- env.SensorMap
	})
}

// Gracefully stops getting data from the broker
func (c * Client) StopGettingData(topic string, rtCh chan map[string]Sensor, quit chan bool) {
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

func (c * Client) SelectEnv(envName string) {

	request := struct {
		Function string `json:"function"`
		EnvName string `json:"data"`
		Auth string `json:"auth"`
	}{
		"selectEnv",
		envName,
		c.LoggedUser,
	}

	jsonReq, _ := json.Marshal(request)

	rx := make(chan []byte)

	//Returns broker's address on rx channel
	go utils.SocketCommunicate(jsonReq, rx)

	// TODO To change test data
	var res = struct {
		Address string `json:"address"`
		Topic string `json:"topic"`
	}{
		"mqtt://user:pass@localhost:1883",
		"testtopic",
	}

	//json.Unmarshal(<- rx, &res)

	uri, err := url.Parse(res.Address)
	if err != nil {
		log.Fatal(err)
	}

	if c.MQTTClient == nil {
		c.MQTTClient = utils.MQTTConnect(c.ID, uri)
	}

}