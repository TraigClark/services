package mqttfile

import (
	"log"
	"main/internal/configpkg"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//copy - pasted from original main.go

func MqttInit(config *configpkg.Config) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		if token.Error() != nil {
			log.Println(token.Error())
		}
	}

	return mqttClient
}

func Publish(client mqtt.Client, topic string, payload string) {
	token := client.Publish(topic, 0, false, payload)
	status := token.WaitTimeout(time.Duration(5) * time.Second)
	if !status {
		log.Println("its cooked")
	}
}
