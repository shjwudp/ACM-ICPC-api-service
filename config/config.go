package config

import (
	"os"

	"encoding/json"
)

// Config is configuration struct
var Config = struct {
	Server struct {
		JWTSecret string
		Port      string
		Admin     struct {
			Account  string
			Password string
		}
	}
	// use sqlite3
	Storage struct {
		Path string
	}
}{}

// LoadConfJSON provide load json config.
func LoadConfJSON(confPath string) (Config, error) {
	var config Config

	configFile, err := os.Open(confPath)

	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if err != nil {
		return config, err
	}

	return config, nil
}
