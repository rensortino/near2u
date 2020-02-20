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

func (c * Client) GetSensorData(uri *url.URL, topic string, rtCh chan map[string]Sensor) {

	client := utils.Connect(c.ID, uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {

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
		/*
			// TODO Define sensor map, {sensorID, measurement}
			 Map needed to update values coming from MQTT broker directly using ID
			 e.g. sensorMap[ID] = sensorData
			 sensorMap = Json.Decode(receivedData)["data"]
			 return sensorList
		*/
		rtCh <- env.SensorMap

		log.Println("Received sensor map: ")
		log.Println(env.SensorMap)
		log.Println("Payload: " + string(msg.Payload()))
	})
}

func (c * Client) SelectEnv(envName string, urlCh chan *url.URL) {

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
	<- rx
	close(rx)


	uri, err := url.Parse(res.Address)
	if err != nil {
		log.Fatal(err)
	}

	urlCh <- uri

}