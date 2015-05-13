package conf

import (
	"encoding/json"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/util"
	"io/ioutil"
)

var mainConfig *Config

type Config struct {
	BaseUrl string
	Proxy   string
	Port    string
	Data    string
	Cors    string
}

func init() {
	mainConfig = new(Config)
}

func Get() *Config {
	return mainConfig
}

func Load(filename string) *Config {
	if filename == "default" {
		filename = "./config.json"
		log.Debug("No config argument given, using default %s", filename)
	}

	if !util.FileExists(filename) {
		log.Critical("Cannot read config file")
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Critical("Cannot load file %s, %s", filename, err)
	}

	mainConfig = new(Config)

	err = json.Unmarshal(data, mainConfig)

	if err != nil {
		log.Critical("Cannot parse config %s", err)
	}

	return mainConfig
}
