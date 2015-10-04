package config

import (
	"io"

	"github.com/naoina/toml"
)

type ServerConfig struct {
	Addr string
}

type StorageConfig struct {
	Driver  string
	Options map[string]interface{}
}

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: ":9009",
		},
		Storage: StorageConfig{
			Driver: "redis",
		},
	}
}

func Parse(configData io.Reader) (*Config, error) {
	conf := DefaultConfig()
	decoder := toml.NewDecoder(configData)

	if err := decoder.Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
