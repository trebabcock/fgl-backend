package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config holds the database parameters
type Config struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

// GetConfig loads the database configuration from the file
func GetConfig() *Config {
	config := Config{}
	configFile, err := os.Open("db/dbConfig.json")
	if err != nil {
		fmt.Println("error opening dbConfig.json:", err)
	}
	byteValue, _ := ioutil.ReadAll(configFile)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		fmt.Println("error unmarshalling dbConfig.json:", err)
	}
	defer configFile.Close()

	return &config
}
