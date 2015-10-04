package redis

import (
	"time"

	"github.com/gamegos/scotty/storage"
	redigo "github.com/garyburd/redigo/redis"
)

func init() {
	storage.Register("redis", initDriver)
}

// Config holds config data for Redis connection.
type Config struct {
	Network     string
	Addr        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Wait        bool
}

// New initializes storage with the given config.
func New(conf *Config) *RedisStorage {
	pool := &redigo.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Wait:        conf.Wait,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial(conf.Network, conf.Addr)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &RedisStorage{pool}
}

func initDriver(config map[string]interface{}) storage.Storage {
	conf, err := configFromMap(config)
	if err != nil {
		panic("adapter:redis: invalid config. " + err.Error())
	}

	return New(conf)
}

func configFromMap(data map[string]interface{}) (*Config, error) {
	conf := &Config{
		Network:     "tcp",
		Addr:        ":6379",
		MaxIdle:     1000,
		MaxActive:   10000,
		IdleTimeout: 60,
	}

	if v, ok := data["network"]; ok {
		conf.Network = v.(string)
	}

	if v, ok := data["addr"]; ok {
		conf.Addr = v.(string)
	}

	if v, ok := data["maxIdle"]; ok {
		conf.MaxIdle = int(v.(int64))
	}

	if v, ok := data["maxActive"]; ok {
		conf.MaxActive = int(v.(int64))
	}

	if v, ok := data["idleTimeout"]; ok {
		conf.IdleTimeout = int(v.(int64))
	}

	if v, ok := data["wait"]; ok {
		conf.Wait = v.(bool)
	}

	return conf, nil
}
