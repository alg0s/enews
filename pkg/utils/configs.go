package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

const configPath = `configs/configs.json`

// Database represents database configs
type Database struct {
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Configs represents configs.json
type Configs struct {
	Database `json:"database"`
}

// LoadConfigs loads user settings from `configs/configs.json` and returns a Configs object
func LoadConfigs() *Configs {
	var cpath, err = filepath.Abs(configPath)

	if err != nil {
		log.Fatal("No path: ", err)
	}

	print("Config path: ", cpath)

	data, err := ioutil.ReadFile(cpath)

	if err != nil {
		log.Fatal("Unable to read config file: ", err)
	}

	var configs Configs
	err = json.Unmarshal(data, &configs)
	if err != nil {
		log.Fatal("Unable to marshal json: ", err)
	}
	return &configs
}
