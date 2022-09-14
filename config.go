package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis_rate/v9"
	"io/ioutil"
	"os"
)

type Config struct {
	ListenPort   int                  `json:"listen_port"`
	RedisAddress string               `json:"redis_address"`
	Namespaces   map[string]Namespace `json:"namespaces"`
}

type Namespace struct {
	Rate   int    `json:"rate"`
	Period string `json:"period"`
}

func (n Namespace) getLimit() redis_rate.Limit {
	switch n.Period {
	case "SECOND":
		return redis_rate.PerSecond(n.Rate)
	case "MINITE":
		return redis_rate.PerMinute(n.Rate)
	case "HOUR":
		return redis_rate.PerHour(n.Rate)
	default:
		panic(fmt.Sprintf("Invalid value for Period:  %s", n.Period))
	}
}

func LoadConfig() (Config, error) {
	filePath := os.Getenv("CONFIG_FILE")
	if len(filePath) == 0 {
		filePath = "/etc/rate_limiter_svc/conf.json"
	}
	configFile, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file %s: %w", filePath, err)
	}
	defer configFile.Close()

	configBytes, _ := ioutil.ReadAll(configFile)
	var config Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config json: %w", err)
	}
	// TODO: validate

	return config, nil
}
