package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port    int
	Timeout int
	// Add more configuration fields here
}

func LoadConfig() (*Config, error) {
	// Set the default values for configuration fields
	viper.SetDefault("port", 8080)
	viper.SetDefault("timeout", 30)

	// Load configuration from a YAML file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./app/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config file: %v", err)
		return nil, err
	}

	// Map the configuration fields to the Config struct
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config file: %v", err)
		return nil, err
	}

	return &cfg, nil
}
