package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Debug        bool `envconfig:"DEBUG" default:"true"`
	DebugVerbose bool `envconfig:"DEBUG_VERBOSE"`
	// MinimalLogLevel is a level for setup minimal logger event notification.
	// Allowed: debug, info, warn, error, dpanic, panic, fatal
	MinimalLogLevel string `envconfig:"MIN_LOG_LEVEL" default:"info"`
	GrpcPort        string `envconfig:"GRPC_PORT" default:"8080"`
	Telegram
	Db
	Proxy
	Metric
}

// Prepare variables to static configuration
func (c *App) Prepare() (err error) {
	if err = envconfig.Process("", c); err != nil {
		return err
	}

	if err = c.Telegram.Prepare(); err != nil {
		return err
	}

	if err = c.Db.Prepare(); err != nil {
		return err
	}

	if err = c.Metric.Prepare(); err != nil {
		return err
	}

	return c.Proxy.Prepare()
}

func (c *App) IsDebug() bool {
	return c.Debug
}

func (c *App) GetMinimalLogLevel() string {
	return c.MinimalLogLevel
}

func (c *App) TelegramCfg() Telegram {
	return c.Telegram
}

func (c *App) GRPCAddress() string {
	return ":" + c.GrpcPort
}
