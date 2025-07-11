package config

import (
	"fmt"
)

type Config struct {
	ServerPort string
	DBConn     string
}

func (c *Config) String() string {
	return fmt.Sprintf("Port: %s, DBConn: %s", c.ServerPort, c.DBConn)
}
