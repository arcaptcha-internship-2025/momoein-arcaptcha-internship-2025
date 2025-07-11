package config

import (
	"encoding/json"
	"os"
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
