package configs

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	defaultPositionRefreshDuration = 10 * time.Second
	defaultTradersRefreshDuration  = 24 * time.Hour
)

type Provider struct {
	PositionRefreshDuration *time.Duration `envconfig:"PROVIDER_POSITION_REFRESH_DURATION"`
	TradersRefreshDuration  *time.Duration `envconfig:"PROVIDER_TRADERS_REFRESH_DURATION"`
}

func (c *Provider) Prepare() error {
	return envconfig.Process("", c)
}

func (c *Provider) PositionRefreshTimeDuration() time.Duration {
	if c.PositionRefreshDuration != nil {
		return *c.PositionRefreshDuration
	}

	return defaultPositionRefreshDuration
}

func (c *Provider) TradersRefreshTimeDuration() time.Duration {
	if c.TradersRefreshDuration != nil {
		return *c.TradersRefreshDuration
	}

	return defaultTradersRefreshDuration
}
