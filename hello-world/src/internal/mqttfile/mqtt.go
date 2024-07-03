package mqttfile

import (
	"log"
	"main/internal/configpkg"
	"time"
	"os"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//copy - pasted from original main.go

func MqttInit(config *configpkg.Config) mqtt.Client {
	broker := os.Getenv("MQTT_BROKER_ADDRESS")
    if broker == "" {
        fmt.Println("MQTT_BROKER_ADDRESS environment variable is not set")
    }

    opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", broker))
    opts.SetClientID("go_mqtt_client")

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Println(token.Error())
    }

    log.Println("Connected to MQTT broker")
	
	return client
    // Your MQTT client logic here
}

func Publish(client mqtt.Client, topic string, payload string) {
	token := client.Publish(topic, 0, false, payload)
	status := token.WaitTimeout(time.Duration(5) * time.Second)
	if !status {
		log.Println("its cooked")
	}
}
