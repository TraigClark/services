package modbusmaker

//imported from external package for modbus and mqtt
import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goburrow/modbus"
)

func ModbusClient(modbus.ClientHandler, modbus.Client) {
	// Modbus TCP
	handler := modbus.NewTCPClientHandler("127.0.0.1:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0
	handler.Logger = log.New(os.Stdout, "modbus: ", log.LstdFlags)

	if err := handler.Connect(); err != nil {
		fmt.Printf("Error connecting to Modbus server: %v\n", err)
	}
	defer handler.Close()
	
	client1 := modbus.NewClient(handler)
	
	// Slave 1
	handler.SlaveId = 0
	client1.WriteSingleRegister(0, 12)
	client1.WriteSingleRegister(321, 52)
	client1.WriteSingleRegister(503, 492)

	client2 := modbus.NewClient(handler)

	// Slave 2
	handler.SlaveId = 1
	client2.WriteSingleRegister(4, 34)
	client2.WriteSingleRegister(322, 54)
	client2.WriteSingleRegister(521, 493)
}