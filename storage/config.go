package storage

import (
	"encoding/json"
	"log"
	"os"
)

// RedisConfig holds config data for Redis connection.
type RedisConfig struct {
	Network     string `json:"network"`
	Addr        string `json:"addr"`
	MaxIdle     int    `json:"max_idle"`
	MaxActive   int    `json:"max_active"`
	IdleTimeout int64  `json:"idle_timeout"`
	Wait        bool   `json:"wait"`
}

// Config holds general config data.
type Config struct {
	Addr  string      `json:"addr"`
	Redis RedisConfig `json:"redis"`
}

// Load imports config data from file into Config struct.
func (conf *Config) Load(confFile *os.File) error {
	decoder := json.NewDecoder(confFile)

	if err := decoder.Decode(conf); err != nil {
		return err
	}

	return nil
}

// InitConfig prepares config file.
func InitConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Can't read config file: '%s'. [error]: %s", configFile, err)
	}

	conf := Config{}
	conf.Redis = RedisConfig{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 60,
		Wait:        true,
	}

	if err := conf.Load(file); err != nil {
		log.Fatal(err)
	}

	return conf
}
