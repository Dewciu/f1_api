package config

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}
	Server struct {
		Host              string `yaml:"host"`
		Port              int    `yaml:"port"`
		Mode              string `yaml:"mode"`
		ApiSecret         string `yaml:"api_secret"`
		TokenHourLifetime int    `yaml:"token_hour_lifetime"`
	}
}

// TODO: Remove this global variable and parse it as an argument
var CONFIG_PATH string

// getDataFromConfig reads the contents of a file specified by file path and returns the data as a byte slice.
// If an error occurs while opening or reading the file, it returns nil and the error.
func GetConfig() (*Config, error) {
	log.Debug("Reading configuration...")
	file, err := os.Open(CONFIG_PATH)
	if err != nil {
		log.Errorf("Failed to open file: %v", err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("Failed to read file: %v", err)
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Errorf("Failed to unmarshal YAML data: %v", err)
		return nil, err
	}
	log.Debug("Database configuration loaded successfully")

	return config, nil
}
