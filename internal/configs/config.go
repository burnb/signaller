package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Debug bool `envconfig:"DEBUG" default:"false"`
	Logger
	Db
	GRPC
	Proxy
	Provider
	Telegram
	Metric
}

// Prepare variables to static configuration
func (c *App) Prepare() (err error) {
	if err = envconfig.Process("", c); err != nil {
		return err
	}

	if err = c.Logger.Prepare(c.Debug); err != nil {
		return err
	}

	if err = c.Db.Prepare(); err != nil {
		return err
	}

	if err = c.GRPC.Prepare(); err != nil {
		return err
	}

	if err = c.Proxy.Prepare(); err != nil {
		return err
	}

	if err = c.Provider.Prepare(); err != nil {
		return err
	}

	if err = c.Telegram.Prepare(); err != nil {
		return err
	}

	return c.Metric.Prepare()
}
