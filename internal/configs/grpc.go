package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type GRPC struct {
	Port string `envconfig:"GRPC_PORT" default:"8080"`
}

// Prepare variables to static configuration
func (c *GRPC) Prepare() error {
	return envconfig.Process("", c)
}

func (c *GRPC) Address() string {
	return ":" + c.Port
}
