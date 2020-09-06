package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

const settingPath = `configs/configs.json`

// Database represents database configs
type Database struct {
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// VnNLPConfigs represents optional configs of VnNLP
type VnNLPConfigs struct {
	StopAfterProgramQuit  bool `json:"stopAfterProgramQuit"`
	NumberArticlePerBatch int  `json:"numberArticlePerbatch"`
}

// VnNLP represents VnNLP server configs
type VnNLP struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Annotators   string `json:"annotators"`
	MaxHeapSize  string `json:"maxHeapSize"`
	VnNLPConfigs `json:"configs"`
}

// Configs represents configs.json
type Configs struct {
	Database `json:"database"`
	VnNLP    `json:"vnnlp"`
}

// LoadConfigs loads user configs from `configs/configs.json` and returns a Configs object
func LoadConfigs() *Configs {
	var cpath, err = filepath.Abs(settingPath)
	if err != nil {
		log.Fatal("No path: ", err)
	}

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
