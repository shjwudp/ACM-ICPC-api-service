package model

import (
	"encoding/json"
	"os"
	"time"
)

// Configuration for ACM-ICPC-api-service
type Configuration struct {
	Server struct {
		JWTSecret string
		Addr      string
		Admin     struct {
			Account  string
			Password string
		}
		IsTestMode bool
	}
	// use sqlite3
	Storage struct {
		Dirver string
		Config string
	}
	Printer struct {
		QueueSize      int
		PinterNameList []string
	}
	ResultsXMLPath string
	ContestInfo    struct {
		StartTime      time.Time
		GoldMedalNum   int
		SilverMedalNum int
		BronzeMedalNum int
		Duration       time.Duration
	}
}

// ConfigurationLoad load Configuration from file
func ConfigurationLoad(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	conf := new(Configuration)
	err = decoder.Decode(conf)
	return conf, err
}
