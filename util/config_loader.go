package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var configInstance *Config

type Config struct {
	GmailID string
	GmailPW string
}

func InitConfig(path string) {
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)

	configInstance = &Config{}
	err := decoder.Decode(configInstance)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func GetConfigInstance() *Config {
	if configInstance == nil {
		panic(errors.New("init config first"))
	}
	return configInstance
}
