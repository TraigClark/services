package service

import (
	"fmt"
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
	handler := modbus.NewTCPClientHandler("127.0.0.1:502")
	client := modbus.NewClient(handler)

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

	// Set the timeout and slave ID for the Modbus handler
	handler.Timeout = time.Duration(config.Timeout) * time.Millisecond
	handler.SlaveId = config.SlaveId

	// Initialize the group size for register grouping
	groupSize := 125

	// Initialize the selected registers to be sampled
	selectedRegister := config.Tags
	// []int{0, 321, 503, 4, 322, 521}

	// Create an empty slice to store the groups of registers
	groups := make([][]int, 0)

	for {
		select {
		case <-stopper:
			return
		case <-sampler.C:
			// map creation
			existingGroups := make(map[int]bool)

			// mark existing groups
			for _, existingGroup := range groups {
				existingGroups[existingGroup[0]] = true
			}

			// iterates over selected registers
			for _, register := range selectedRegister {
				registerValue := register.Register

				group := int(registerValue) / groupSize

				// check if already exists
				if !existingGroups[group] {
					groups = append(groups, []int{group})
					existingGroups[group] = true
				}
				
				if connected {
					// sampling the register
					for _, group := range groups {
						values, err := client.ReadHoldingRegisters(uint16(group[0]*groupSize), uint16(groupSize))
						if err != nil {
							// err handling
							fmt.Println("Error sampling int:", err)
							connected = false
						} else {
							for i, value := range values {
								if value != 0 {
									fmt.Printf("Sampled register %d from slave %d with value %d\n", registerValue, config.SlaveId, value)
									mqttfile.Publish(ret, config.Mqtttopic, fmt.Sprintf("%v", value))
									fmt.Println(values[i])
								}
							}
						}
					}
				} else {
					fmt.Println("Not connected, skipping sampling")

					var err error
					fmt.Println("Error connecting:", err)
					connected = false
				}
			}
		case <-reconnectTicker.C:
			if !connected {
				handler.Close()
				if err := handler.Connect(); err != nil {
					fmt.Println("Reconnection Failed")
					fmt.Println(err)
				} else {
					client = modbus.NewClient(handler)
					fmt.Println("Reconnection worked")
					connected = true
				}
			}
		}
	}
}

// Helper function to check if a group already exists in the groups array
