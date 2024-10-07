package common

import (
	"encoding/json"
	"log"
	"os"
)

var Config Configuration

type Configuration struct {
	ConnectionString string
	ServerIp         string
	ServerPort       string
}

func LoadConfig() Configuration {

	file, err := os.Open("config/config.json")

	if err != nil {
		log.Fatalln("cannot open configuration file", err)
	}

	decoder := json.NewDecoder(file)

	Config = Configuration{}

	err = decoder.Decode(&Config)

	if err != nil {
		log.Fatalln("error occured in decoding configuration content", err)
	}

	file.Close()

	return Config
}

func PrintJsonFormat(payload string, req interface{}) {

	result, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		log.Println("Error while Json marshal", err.Error())
	}

	log.Println(payload+" :  ", string(result))
}
