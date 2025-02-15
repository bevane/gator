package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	config := Config{}
	configFilePath := getConfigFilePath()
	data, _ := os.ReadFile(configFilePath)
	json.Unmarshal(data, &config)
	return config
}

func (c *Config) SetUser(username string) {
	configFilePath := getConfigFilePath()
	c.CurrentUserName = username
	data, _ := json.Marshal(c)
	os.WriteFile(configFilePath, data, 0666)
}

func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir()
	configFilePath := filepath.Join(homeDir, "gatorconfig.json")
	return configFilePath
}
