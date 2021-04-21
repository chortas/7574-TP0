package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

type ClientConfigJSON struct {
	ServerAddress string `json:"CLI_SERVER_ADDRESS"`
	ID            string `json:"CLI_ID"`
	LoopLapse     string `json:"CLI_LOOP_LAPSE"`
	LoopPeriod    string `json:"CLI_LOOP_PERIOD"`
}

func LoadConfigFromFile() (common.ClientConfig) {
	jsonFile, err := os.Open("./config/client_config.json")
	if err != nil {
		log.Fatalf("%s", err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var clientConfig common.ClientConfig
	var clientConfigJSON ClientConfigJSON
	err = json.Unmarshal(byteValue, &clientConfigJSON)
	if err != nil {
		log.Fatalf("%s", err)
	}
	loop_lapse, err := time.ParseDuration(clientConfigJSON.LoopLapse)
	if err != nil {
		log.Fatalf("%s", err)
	}
	loop_period, err := time.ParseDuration(clientConfigJSON.LoopPeriod)
	if err != nil {
		log.Fatalf("%s", err)
	}
	clientConfig = common.ClientConfig{
		ServerAddress: clientConfigJSON.ServerAddress,
		ID:            clientConfigJSON.ID,
		LoopLapse:     loop_lapse,
		LoopPeriod:    loop_period,
	}
	return clientConfig
}

// InitConfig Function that uses viper library to parse env variables. If
// some of the variables cannot be parsed, an error is returned
func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper to read env variables with the CLI_ prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")

	// Add env variables supported
	v.BindEnv("id")
	v.BindEnv("server", "address")
	v.BindEnv("loop", "period")
	v.BindEnv("loop", "lapse")

	// Parse time.Duration variables and return an error
	// if those variables cannot be parsed
	if _, err := time.ParseDuration(v.GetString("loop_lapse")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_LAPSE env var as time.Duration.")
	}

	if _, err := time.ParseDuration(v.GetString("loop_period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
	}

	return v, nil
}

func LoadConfigFromEnvVariables() (common.ClientConfig) {
	v, err := InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	clientConfig := common.ClientConfig{
		ServerAddress: v.GetString("server_address"),
		ID:            v.GetString("id"),
		LoopLapse:     v.GetDuration("loop_lapse"),
		LoopPeriod:    v.GetDuration("loop_period"),
	}

	return clientConfig
}


func main() {
	clientConfig := LoadConfigFromFile()
	log.Printf("Client config: %v", clientConfig)


	//clientConfig := LoadConfigFromEnvVariables()

	client := common.NewClient(clientConfig)
	client.StartClientLoop()
}
