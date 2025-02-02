package main

import (
	//imported from internal package
	"main/internal/configpkg"
	"main/internal/mqttfile"
	"main/internal/service"
	"main/internal/modbusmaker"

	//imported from external package
	"log"
	"sync"
	"time"
	"github.com/goburrow/modbus"
)

func main() {
	// Read JSON file
	config1, err := configpkg.ReadConfigFromFile("config/config.json")
	if err != nil {
		log.Printf("Error reading config file: %v\n", err)
		return
	}

	// mqttInit
	ret := mqttfile.MqttInit(config1)

	var wg sync.WaitGroup
	stopper := make(chan struct{})

	for _, device := range config1.Devices {
		wg.Add(1)
		go service.GoRoutine(&wg, stopper, &device, ret)
	}

	time.Sleep(100 * time.Second)
	close(stopper)
	wg.Wait()
	
	for _, device := range config1.Devices {
		log.Println(device)
		handler := modbus.NewTCPClientHandler("host.docker.internal:1502")
		client := modbus.NewClient(handler)
		modbusmaker.ModbusClient(handler, client)
	}
}
