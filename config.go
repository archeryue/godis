package main

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func LoadConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	jsonStr, err := io.ReadAll(file)
	if err != nil {
		return
	}

	config = &Config{}
	if err = json.Unmarshal(jsonStr, config); err != nil {
		return nil, err
	}
	return
}
