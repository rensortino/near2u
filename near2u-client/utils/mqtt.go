package utils

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"time"
)

type Sensor struct {
	ID string `json:ID`
	Name string `json:name`
	Measurement float32 `json:measurement`
}

type Environment struct {
	ID string `json:ID`
	SensorMap map[string]Sensor `json:sensors`
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func Listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {

		sensorList := make(map[string]Sensor)
		env := Environment{SensorMap:sensorList}

		/*
			 Payload format:
			{"ID":"env1","SensorMap":{
				"sensor1":{"ID":"id","Name":"name","Measurement":7.4 },
				"sensor2": {"ID":"otherID","Name":"name2","Measurement":4.76}
			}}
		*/
		json.Unmarshal(msg.Payload(), &env)
		fmt.Println("Received data: ")
		fmt.Println(env)
		fmt.Println("Payload: " + string(msg.Payload()))

	})
}