package client

import (
	"../utils"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Client struct {
	ID         string
	LoggedUser string
	MQTTClient mqtt.Client
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

// Used when configuring an environment, to indicate which one to work with
func SetCurrentEnv(currentEnv * Environment, name string) {
	if currentEnv.Name != name {
		currentEnv.Name = name
		currentEnv.DeviceMap = make(map[string]interface{}) // Initialize the map to populate it with the devices
		currentEnv.LastModified = 0
	}
}

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func (c *Client) GetEnvList(envNameCh chan string, errCh chan string) {

	res := utils.SocketCommunicate("visualizza_ambienti", c.LoggedUser, nil)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	envNameCh <- "start" // Used to enter the correct case in the select
	for _ ,v := range res["data"].(map[string]interface{})["environments"].([]interface{}) {
		envNameCh <- v.(string)
	}
	close(envNameCh)
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

	//Returns broker's address
	res := utils.SocketCommunicate("topic_ambiente", c.LoggedUser, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	topicCh <- res["data"].(map[string]interface{})["topic"].(string)
	uriCh <- res["data"].(map[string]interface{})["broker_host"].(string)
	close(topicCh)
	close(uriCh)
}
/*
func (c *Client) SelectEnv(envName string, envCh chan *Environment, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	//Returns broker's address on rx channel
	res := utils.SocketCommunicate("seleziona_ambiente", c.LoggedUser, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}

	envCh <- res["data"].(map[string]interface{})["environment"].(*Environment)
	close(envCh)
}
*/
func (c *Client) CreateEnv(envName string, currentEnv * Environment, resCh, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	res := utils.SocketCommunicate("crea_ambiente", c.LoggedUser, data)

	if res["status"] == "Succesfull" {
		SetCurrentEnv(currentEnv, envName)
		resCh <- res["status"].(string)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
}
