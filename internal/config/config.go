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

func Read() (Config, error) {
	config := Config{}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return config, err
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	c.CurrentUserName = username
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(homeDir, ".gatorconfig.json")
	return configFilePath, nil
}
