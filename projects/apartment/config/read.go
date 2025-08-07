package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v11"
)

func MustReadJson(path string) Config {
	cfg, err := ReadJson(path)
	if err != nil {
		panic(err)
	}
	return cfg
}

func ReadJson(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	return cfg, json.Unmarshal(b, &cfg)
}

func ReadEnv() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func MustReadEnv() Config {
	cfg, err := ReadEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}
