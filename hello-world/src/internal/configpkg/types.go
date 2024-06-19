package configpkg

type Config struct {
	Devices []DeviceConfig `json:"devices"`
}

type DeviceConfig struct {
	Mqttbroker    string      `json:"mqttbroker"`
	Mqtttopic     string      `json:"mqtttopic"`
	SlaveId       uint8       `json:"slave_id"`
	Timeout       int         `json:"timeout"`
	SampleRate    int         `json:"sample_rate"`
	ReconnectRate int         `json:"reconnect_rate"`
	Tags          []TagConfig `json:"tags"`
}

type TagConfig struct {
	Register     uint16  `json:"register"`
	RegisterType string  `json:"register_type"`
	TagName      string  `json:"tag_name"`
	Multiplier   float64 `json:"multiplier"` // x
	Offset       float64 `json:"offset"` // b
	// slope intercept form y = mx + b 
	// if both are 0 do nothing
	//make file, docker file, containerizing
}