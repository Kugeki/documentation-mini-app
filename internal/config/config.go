package config

import (
	"encoding/json"
	"io"
	"log"
)

type Config struct {
	Addr string `json:"addr"`
}

func Parse(r io.Reader) *Config {
	var config Config
	err := json.NewDecoder(r).Decode(&config)
	if err != nil {
		log.Fatalf("config json decode: %v\n", err)
	}

	return &config
}
