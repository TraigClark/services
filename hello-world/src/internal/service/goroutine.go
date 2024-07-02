package service

import (
	"encoding/binary"
	"fmt"
	"log"
	"main/internal/configpkg"
	"main/internal/modbusmaker"
	"main/internal/mqttfile"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/modbus"
)

func GoRoutine(wg *sync.WaitGroup, stopper chan struct{}, config *configpkg.DeviceConfig, ret mqtt.Client) {
	// Create a new TCP client handler for Modbus communication
	handler := modbus.NewTCPClientHandler("host.docker.internal:1502")
	client := modbus.NewClient(handler)
	log.Println("handler and client initalized")

	// Initialize the Modbus client
	modbusmaker.ModbusClient(handler, client)

	defer wg.Done()

	var connected bool

	// Create a ticker for periodic sampling
	sampler := time.NewTicker(time.Duration(config.SampleRate) * time.Second)
	defer sampler.Stop()

	// Create a ticker for periodic reconnection
	reconnectTicker := time.NewTicker(time.Duration(config.ReconnectRate) * time.Second)
	defer reconnectTicker.Stop()

	mqttTicker := time.NewTicker(time.Duration(config.ReconnectRate) * time.Second)
	defer mqttTicker.Stop()

	log.Println("tickers initalized")

	// Set the timeout and slave ID for the Modbus handler
	//change back
	handler.Timeout = time.Duration(100) * time.Second
	handler.SlaveId = config.SlaveId

	// Initialize the group size for register grouping
	groupSize := 125

	for {
		select {
		case <-stopper:
			return
		case <-sampler.C:
			if connected {
				log.Println("connected")
				//ware := config.Tags[0].TagName
				groups, registers := modbusmaker.OrganizeRegisters(config)
				// sampling the register
				for j := 0; j < len(groups); j++ {
					values, err := client.ReadHoldingRegisters(uint16(groups[j][0]*groupSize), uint16(groupSize))
					// valuesInput, err := client.ReadInputRegisters(uint16(groups[j][0]*groupSize), uint16(groupSize))
					if err != nil {
						// error handling
						log.Println("Error sampling input registers:", err)
						connected = false
					} else {
						// iterate through length of register array
						for i := 0; i < len(registers); i++ {
							id := config.SlaveId
							if registers[i]/groupSize == groups[j][0] {
								value := binary.BigEndian.Uint16(values[((registers[i])-(groupSize*groups[j][0]))*2:])
								// valuesInput := binary.BigEndian.Uint16(valuesInput[((registers[i])-(groupSize*groups[j][0]))*2:])

								if value != 0 {
									log.Printf("Sampled holding register %v from slave %v with value %v\n", registers[i], id, (float64(value)*config.Tags[i].Multiplier)+config.Tags[i].Offset)
									test := (float64(value) * config.Tags[i].Multiplier) + config.Tags[i].Offset
									val1 := config.Tags[i].TagName
									log.Println(fmt.Sprintf("%s: %f", val1, test))
									mqttfile.Publish(ret, config.Mqtttopic, fmt.Sprintf("%s: %f", val1, test))
								}

								// if valueInput != 0 {
								// 	log.Printf("Sampled holding register %v from slave %v with value %v\n", registers[i], config.SlaveId, (float64(valueInput)*config.Tags[i].Multiplier)+config.Tags[i].Offset)
								// 	test := (float64(valueInput) * config.Tags[i].Multiplier) + config.Tags[i].Offset
								// 	val2 := config.Tags[i].TagName
								// 	mqttfile.Publish(ret, config.Mqtttopic, fmt.Sprintf("%s: %f", val2, test))
								// }
							}
						}
					}
				}
			} else {
				log.Println("Not connected, skipping sampling")
				var err error
				log.Println("Error connecting:", err)
				connected = false
			}
		case <-reconnectTicker.C:
			if !connected {
				handler.Close()
				if err := handler.Connect(); err != nil {
					log.Println("Reconnection Failed")
					log.Println(err)
				} else {
					client = modbus.NewClient(handler)
					log.Println("Reconnection worked")
					connected = true
				}
			}
		case <-mqttTicker.C:
			if ret.IsConnected() {
				log.Println("MQTT connection is already established")
			} else {
				log.Println("MQTT connection is not established, reconnecting...")
				mqttfile.MqttInit(&configpkg.Config{})
				if token := ret.Connect(); token.Wait() && token.Error() != nil {
					log.Println("MQTT reconnection failed:", token.Error())
				} else {
					log.Println("MQTT reconnection successful")
				}
			}
		}
	}
}
