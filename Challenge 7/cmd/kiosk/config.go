package main

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	Config Config
}

func (c *AppConfig) LoadConfigurations() (string, error) {
	// check for config.json
	_, err := os.Stat("./conf/config.json")
	if err != nil {
		return "File not found", err
	}
	// read config.json
	data, err := os.ReadFile("./conf/config.json")
	if err != nil {
		return "Failed to read file", err
	}

	// convert to struct
	err = json.Unmarshal(data, &c.Config)
	if err != nil {
		return "Failed to convert to config", err
	}
	// return
	return "", nil
}
