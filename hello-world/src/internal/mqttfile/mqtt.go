package mqttfile

import (
	"github.com/eclipse/paho.mqtt.golang"
	"main/internal/configpkg"
	"log"
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
	token.Wait()
}