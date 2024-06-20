package modbusmaker

//imported from external package for modbus and mqtt
import (
	"log"
	"os"
	"time"
	"main/internal/configpkg"
	"github.com/goburrow/modbus"
)

func ModbusClient(modbus.ClientHandler, modbus.Client) {
	// Modbus TCP
	handler := modbus.NewTCPClientHandler("host.docker.internal:1502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0
	handler.Logger = log.New(os.Stdout, "modbus: ", log.LstdFlags)

	if err := handler.Connect(); err != nil {
		log.Printf("Error connecting to Modbus server: %v\n", err)
	} else {
		log.Println("No error connecting to modbus server")
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	// Slave 1
	handler.SlaveId = 0
	client.WriteSingleRegister(0, 12)
	client.WriteSingleRegister(321, 52)
	client.WriteSingleRegister(503, 492)

	// Slave 2
	handler.SlaveId = 1
	client.WriteSingleRegister(4, 34)
	client.WriteSingleRegister(322, 54)
	client.WriteSingleRegister(521, 493)
}

func OrganizeRegisters(config *configpkg.DeviceConfig)(groupfin [][]int, registersfin []int) {
	// Create an empty slice to store the groups of registers
	groups := make([][]int, 0)
	var registers []int

	// map creation
	existingGroups := make(map[int]bool)

	// Initialize the selected registers to be sampled
	selectedRegister := config.Tags
	// []int{0, 321, 503, 4, 322, 521}
	
	// Initialize the group size for register grouping
	groupSize := newFunction()

	///////////////////////////////////////////////

	// iterates over selected registers
	for _, register := range selectedRegister {
		registerValue := register.Register
		registers = append(registers, int(registerValue))
		
		group := int(registerValue) / groupSize

		// check if group already exists
		if !existingGroups[group] {
			groups = append(groups, []int{group})
			existingGroups[group] = true
		}
	}
	return groups, registers
}

func newFunction() int {
	groupSize := 125
	return groupSize
}
