package config

import (
	"fmt"
)

type Config struct {
	ServerPort string
}

func (c *Config) String() string {
	return fmt.Sprintf("Port: %s", c.ServerPort)
}
