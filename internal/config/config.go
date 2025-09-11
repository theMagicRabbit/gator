package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configJsonName string = ".gatorconfig.json"

type Config struct {
	Db_url string;
	Current_user_name string;
}

// generateConfigFilePath generates the full path name for the config file
// based on the users home directory path
func generateConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFileName := fmt.Sprintf("%s/%s", home, configJsonName)
	return configFileName, nil
}

// Read reads the config file and returns the content as a Config struct
func Read() (Config, error) {
	config := Config{}
	configFileName, err := generateConfigFilePath()
	if err != nil {
		return config, err
	}
	configData, err := os.ReadFile(configFileName)
	if err != nil {
		return config, err
	}
	json.Unmarshal(configData, &config)
	return config, nil
}

// SetUser sets the current user value and writes current configuration to
// the config file.
func(c Config) SetUser(current_user string) error {
	c.Current_user_name = current_user
	configData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	configFileName, err := generateConfigFilePath()
	err = os.WriteFile(configFileName, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}
