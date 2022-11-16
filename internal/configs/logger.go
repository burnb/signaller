package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Logger struct {
	// LogLevel is a level for setup minimal logger event notification.
	// Allowed: debug, info, warn, error, dpanic, panic, fatal
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func (c *Logger) MinimalLevel() string {
	return c.LogLevel
}

func (c *Logger) Prepare(debug bool) error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}

	if debug {
		c.LogLevel = "debug"
	}

	return nil
}
