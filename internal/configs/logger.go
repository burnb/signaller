package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Logger struct {
	// MinimalLevel is a level for setup minimal logger event notification.
	// Allowed: debug, info, warn, error, dpanic, panic, fatal
	MinimalLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func (c *Logger) Prepare(debug bool) error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}

	if debug {
		c.MinimalLevel = "debug"
	}

	return nil
}
