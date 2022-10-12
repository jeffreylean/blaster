package config

import (
	"context"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kelindar/loader"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Schedule  []LoadSchedule `yaml:"schedule"`
	Iteration int            `yaml:"iteration"`
}

type LoadSchedule struct {
	Hour     int `yaml:"hour"`
	Requests int `yaml:"requests"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if u, ok := os.LookupEnv("BLASTER_CONF"); ok {
		l := loader.New()
		b, err := l.Load(context.Background(), u)
		if err != nil {
			return nil, fmt.Errorf("config: unable to load config, due to %s", err.Error())
		}

		// Parse YAML.
		if err := yaml.Unmarshal(b, cfg); err != nil {
			return nil, fmt.Errorf("config: unable to parse yaml, due to %s", err.Error())
		}
	}

	// Fill with environment variables
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("config: unable to parse config, due to %s", err.Error())
	}
	return cfg, nil
}

// LoadOrPanic loads the configuration or panics
func LoadOrPanic() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}

	return cfg
}
