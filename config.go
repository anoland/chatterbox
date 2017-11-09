package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	NodeAddress string `json:"nodeaddress"`
	NodePort    int    `json:"nodeport"`
	ListenPort  int    `json:"listenport"`
}

func main() {

	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Could not open config file. Please create from the sample file")
	}
	config := Config{}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Println("problem parsing config file", err)
	}

	fmt.Printf("%#v", config)
}
