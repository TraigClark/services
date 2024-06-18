package configpkg

import (
	"encoding/json"
	"io"
	"os"
)


func ReadConfigFromFile(filename string) (*Config, error) {
	jSonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jSonFile.Close()

	// Read the contents of the config file
	data, err := io.ReadAll(jSonFile)
	if err != nil {
		return nil, err
	}

	// Decode the JSON data into the Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}
