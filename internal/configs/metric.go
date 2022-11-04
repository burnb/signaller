package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Metric struct {
	Debug    bool
	HttpPort string `envconfig:"METRIC_HTTP_PORT" default:"8000"`
	Path     string `envconfig:"METRIC_PATH" default:"/"`
}

// Prepare variables to static configuration
func (c *Metric) Prepare() error {
	return envconfig.Process("", c)
}

func (c *Metric) Address() string {
	return ":" + c.HttpPort
}

func (c *Metric) HttpPath() string {
	return c.Path
}
