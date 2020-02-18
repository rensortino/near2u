package client

import (
	"../utils"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"strconv"
)

// TODO Create Client struct to contain session data and client data (See TODO on socket.go)

type Sensor struct {
	ID string `json:ID`
	Name string `json:name`
	Measurement float32 `json:measurement`
}

type Environment struct {
	ID string `json:ID`
	SensorMap map[string]Sensor `json:sensors`
}

// TODO Add Client interface (?)

// Gets an array of Environment IDs from the server, to be displayed on the GUI for selection
func GetEnvList(c * utils.Client) []string {

	/*
	request := utils.RequestParams {
		"getEnvList",
		c.Token,
	}

	jsonReq, _ := json.Marshal(request)

	c.SocketSend(conn, jsonReq)

	rx := make(chan []byte)
	c.SocketReceive(rx)

	// TODO Take environment list accessing the json
	var envList struct {
		Environments [] string `json:"environments"`
	}
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

func getSensorData(c * utils.Client, uri *url.URL, topic string) {

	client := utils.Connect(c.ClientID, uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {

		// TODO Substitute with map[string]interface{}
		sensors := make(map[string]Sensor)
		env := Environment{SensorMap:sensors}

		/*
			 Payload format:
			{"ID":"env1","SensorMap":{
				"sensor1":{"ID":"id","Name":"name","Measurement":7.4 },
				"sensor2":{"ID":"otherID","Name":"name2","Measurement":4.76}
			}}
		*/
		json.Unmarshal(msg.Payload(), &env)
		log.Println("Received sensor map: ")
		log.Println(env.SensorMap)
		log.Println("Payload: " + string(msg.Payload()))
	})
}

func SelectEnv(c * utils.Client, rx chan []byte, envName string) {

	params := utils.RequestParams {
		Function: "selectEnv",
		Auth:     c.Token,
	}

	request := struct {
		Params utils.RequestParams `json:"params"`
		EnvName string `json:"envName"`
	}{
		params,
		envName,
	}

	jsonReq, _ := json.Marshal(request)

	c.SocketSend(jsonReq)

	//Returns broker's address
	c.SocketReceive(rx)

	var res = struct {
		Address string `json:"address"`
		Topic string `json:"topic"`
	}{}

	json.Unmarshal(<- rx, &res)

	uri, err := url.Parse(res.Address)
	if err != nil {
		log.Fatal(err)
	}

	go getSensorData(c, uri, res.Topic)


	/*
	// TODO Define sensor map, {sensorID, measurement} (?)
	 Map needed to update values coming from MQTT broker directly using ID
	 e.g. sensorMap[ID] = sensorData
	 sensorMap = Json.Decode(receivedData)["data"]
	 return sensorList
	 */

}