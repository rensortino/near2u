package client

import (
	"../utils"
	"encoding/json"
 	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	ID         string
	LoggedUser * utils.User
	MQTTClient mqtt.Client
}

var clientInstance *Client

// Implements singleton pattern
func GetClientInstance() *Client {

	if clientInstance == nil {
		clientInstance = &Client{
			"ID1",
			&utils.User{},
			nil,
		}
	}

	return clientInstance
}

func NewEnvironment() * Environment{
	return &Environment{}
}

// Used when configuring an environment, to indicate which one to work with
func SetCurrentEnv(currentEnv * Environment, name string) {

	if currentEnv.Name != name {
		currentEnv.Name = name
		// Initialize the maps to populate them with the devices
		currentEnv.SensorMap = make(map[string]Sensor)
		currentEnv.ActuatorMap = make(map[string]Actuator)
		currentEnv.LastModified = 0
	}
}

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func (c *Client) GetEnvList(envNameCh chan string, errCh chan string) {

	res := utils.SocketCommunicate("visualizza_ambienti", c.LoggedUser.Auth, nil)

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

func (c *Client) GetSensorData(topic string, rtCh chan interface{}, startCh chan bool) {

	<- startCh
	c.MQTTClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Executes every time a message is published on the topic
		receivedData := struct {
			Code int
			Name string
			Kind string
			Measurement float64
			Timestamp string
		} {}

		 json.Unmarshal(msg.Payload(), &receivedData)

		rtCh <- &Sensor{
			Device{
				receivedData.Code,
				receivedData.Name,
				receivedData.Kind,
			},
			receivedData.Measurement,
		}
	})
}

// Gracefully stops getting data from the broker
func (c *Client) StopGettingData(topic string, rtCh chan interface{}, quit chan bool) {
	c.MQTTClient.Unsubscribe(topic)
	c.MQTTClient.Disconnect(1000) // Waits for 1 second before disconnecting

	quit <- true
	close(rtCh)
	close(quit)
}

func (c *Client) GetTopicAndUri(envName string, topicCh, uriCh, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	//Returns broker's address
	res := utils.SocketCommunicate("topic_ambiente", c.LoggedUser.Auth, data)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

	topicCh <- res["data"].(map[string]interface{})["topic"].(string)
	uriCh <- res["data"].(map[string]interface{})["broker_host"].(string)
	close(topicCh)
	close(uriCh)
}

func (c *Client) CreateEnv(envName string, currentEnv * Environment, resCh, errCh chan string) {

	data := struct {
		Name string `json:"name"`
	}{
		envName,
	}

	res := utils.SocketCommunicate("crea_ambiente", c.LoggedUser.Auth, data)

	if res["status"] == "Successful" {
		SetCurrentEnv(currentEnv, envName)
		resCh <- res["status"].(string)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
}

func (c *Client) DeleteEnv(envName string, currentEnv * Environment, resCh, errCh chan string) {

	data := struct {
		Name string `json:"envname"`
	}{
		envName,
	}

	res := utils.SocketCommunicate("elimina_ambiente", c.LoggedUser.Auth, data)

	if res["status"] == "Successful" {
		if currentEnv.Name == envName {
			currentEnv = NewEnvironment()
		}
		resCh <- res["status"].(string)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
}
