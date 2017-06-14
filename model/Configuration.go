package model

import (
	"encoding/json"
	"os"
)

// StorageConfiguration about db
type StorageConfiguration struct {
	Dirver       string
	Addr         string
	MaxIdleConns int
	MaxOpenConns int
}

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
		NeedAuth   bool
	}
	// use sqlite3
	Storage StorageConfiguration
	Printer struct {
		QueueSize      int
		PinterNameList []string
	}
	ResultsXMLPath string
	ContestInfo    ContestInfo
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
