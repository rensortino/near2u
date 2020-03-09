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

<<<<<<< HEAD
type Sensor struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Environment struct {
	Name      string                 `json:"envname"`
	SensorMap map[int]interface{} `json:"sensors"`
}

=======
>>>>>>> Iterazione_3
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
<<<<<<< HEAD

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
=======
}

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func (c *Client) GetEnvList(envNameCh chan string, errCh chan string) {
>>>>>>> Iterazione_3


	res := utils.SocketCommunicate("visualizza_ambienti", c.LoggedUser.Auth, nil)

	if res["status"] == "Failed" {
		errCh <- res["error"].(string)
		return
	}

<<<<<<< HEAD
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
=======
	envNameCh <- "start" // Used to enter the correct case in the select
	for _ ,v := range res["data"].(map[string]interface{})["environments"].([]interface{}) {
		envNameCh <- v.(string)
	}
	close(envNameCh)
>>>>>>> Iterazione_3
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
		Name string `json:"envcode"`
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

<<<<<<< HEAD
// TODO Change sensorlist type to []Sensor
func (c * Client) SelectEnv(envName string, envCh chan * Environment, errCh chan string) {
=======
func (c *Client) CreateEnv(envName string, currentEnv * Environment, resCh, errCh chan string) {
>>>>>>> Iterazione_3

	sensorListCh := make(chan [] Sensor)

<<<<<<< HEAD
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
=======
	res := utils.SocketCommunicate("crea_ambiente", c.LoggedUser.Auth, data)

	if res["status"] == "Successful" {
		SetCurrentEnv(currentEnv, res["data"].(map[string]interface{})["code"].(string))
		resCh <- res["status"].(string)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
>>>>>>> Iterazione_3
}

func (c *Client) DeleteEnv(envName string, currentEnv * Environment, resCh, errCh chan string) {

	data := struct {
		Name string `json:"envcode"`
	}{
		envName,
	}

	res := utils.SocketCommunicate("elimina_ambiente", c.LoggedUser.Auth, data)

<<<<<<< HEAD
	res := <-rx

	if res["status"] == "Succesfull" {
		newEnv := &Environment{
			envName,
			make(map[int]interface{}),
=======
	if res["status"] == "Successful" {
		if currentEnv.Name == envName {
			currentEnv = NewEnvironment()
>>>>>>> Iterazione_3
		}
		resCh <- res["status"].(string)
		return
	} else {
		errCh <- res["error"].(string)
<<<<<<< HEAD
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
=======
		close(errCh)
		return
	}
}

// Used by admin to associate a user to an existing environment
func (c * Client) AssociateUser(envName, email string, resCh, errCh chan string) {

	data := struct {
		Name string `json:"envcode"`
		User string `json:"user"`
>>>>>>> Iterazione_3
	}{
		envName,
		email,
	}

<<<<<<< HEAD
	rx := make(chan map[string]interface{})

	go utils.SocketCommunicate(function, c.LoggedUser, data, rx)

	res := <-rx
=======
	res := utils.SocketCommunicate("associa_utente", c.LoggedUser.Auth, data)
>>>>>>> Iterazione_3

	if res["status"] == "Successful" {
		resCh <- res["status"].(string)
		return
	} else {
		errCh <- res["error"].(string)
		close(errCh)
		return
	}
}