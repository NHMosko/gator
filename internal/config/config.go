package config

import (
	"encoding/json"
	"os"
)

const configFilePath = "/.gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName 
	jsonData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir() 
	if err != nil {
		return err
	}

	if err := os.WriteFile(home  + configFilePath, jsonData, 0o600); err != nil {
		return err
	}
	
	return nil
}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	var config Config
	jsonData, err := os.ReadFile(home + configFilePath)
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
