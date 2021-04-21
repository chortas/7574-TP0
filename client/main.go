package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

import (
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

type ClientConfigJSON struct {
	ServerAddress string `json:"CLI_SERVER_ADDRESS"`
	ID            string `json:"CLI_ID"`
	LoopLapse     string `json:"CLI_LOOP_LAPSE"`
	LoopPeriod    string `json:"CLI_LOOP_PERIOD"`
}

func InitConfigFromFile() (ClientConfigJSON, error) {
	var clientConfigJSON ClientConfigJSON

	jsonFile, err := os.Open("./config/client_config.json")
	if err != nil {
		return clientConfigJSON, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &clientConfigJSON)
	if err != nil {
		return clientConfigJSON, err
	}

	return clientConfigJSON, nil
}

// InitConfigFromEnvVariables Function that uses viper library to parse env variables. If
// some of the variables cannot be parsed, an error is returned
func InitConfigFromEnvVariables() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper to read env variables with the CLI_ prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")

	// Add env variables supported
	v.BindEnv("id")
	v.BindEnv("server", "address")
	v.BindEnv("loop", "period")
	v.BindEnv("loop", "lapse")

	return v, nil
}

func LoadConfig() (common.ClientConfig) {
	vEnv, errEnv := InitConfigFromEnvVariables()
	LogError(errEnv)

	vConfig, errConfig := InitConfigFromFile()
	LogError(errConfig)

	var server_address string 
	var id string            
	var loop_lapse time.Duration     
	var loop_period time.Duration
	var err error

	if vEnv.IsSet("server_address") {
		server_address = vEnv.GetString("server_address")
	} else if vConfig.ServerAddress != "" {
		server_address = vConfig.ServerAddress
	} else {
		server_address = "server:12345"
	}

	if vEnv.IsSet("id") {
		id = vEnv.GetString("id")
	} else if vConfig.ID != "" {
		id = vConfig.ID
	} else {
		id = "1"
	}

	if vEnv.IsSet("loop_lapse") {
		loop_lapse, err = time.ParseDuration(vEnv.GetString("loop_lapse"))
		LogError(err)
	} else if vConfig.LoopLapse != "" {
		loop_lapse, err = time.ParseDuration(vConfig.LoopLapse)
		LogError(err)
	} else {
		loop_lapse, err = time.ParseDuration("1m2s")
		LogError(err)
	}

	if vEnv.IsSet("loop_period") {
		loop_period, err = time.ParseDuration(vEnv.GetString("loop_period"))
		LogError(err)
	} else if vConfig.LoopPeriod != "" {
		loop_period, err = time.ParseDuration(vConfig.LoopPeriod)
		LogError(err)
	} else {
		loop_period, err = time.ParseDuration("10s")
		LogError(err)
	}

	clientConfig := common.ClientConfig{
		ServerAddress: server_address,
		ID:            id,
		LoopLapse:     loop_lapse,
		LoopPeriod:    loop_period,
	}

	return clientConfig
}

func LogError(err error) () {
	if (err != nil) {
		log.Fatalf("%s", err)
	}
}

func main() {
	clientConfig := LoadConfig()
	log.Printf("Client config: %v", clientConfig)

	client := common.NewClient(clientConfig)
	client.StartClientLoop()
}
