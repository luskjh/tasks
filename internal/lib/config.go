package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port uint16
}

func NewConfig(path string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
