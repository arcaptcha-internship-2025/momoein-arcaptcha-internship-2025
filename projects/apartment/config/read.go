package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
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

func ReadEnv(filenames ...string) (Config, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return Config{}, nil
	}
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func MustReadEnv(filenames ...string) Config {
	cfg, err := ReadEnv(filenames...)
	if err != nil {
		panic(err)
	}
	return cfg
}
