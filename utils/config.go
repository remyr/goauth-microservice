package utils

import (
	"os"
	"fmt"
	"encoding/json"
)

type Database struct {
	Host 		string			`json:"host"`
	Name 		string			`json:"name"`
}

type Config struct {
	Database 	Database 		`json:"database"`
	Port 		string 			`json:"port"`
}

func LoadConfiguration(file string) Config {
	config := Config{}
	if file == "" {
		file = "config.json"
	}
	configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)

	return config
}