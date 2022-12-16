package configs

import (
	"math/rand"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	defaultPositionRefreshDuration = 10 * time.Second
	defaultTradersRefreshDuration  = 24 * time.Hour
)

type Provider struct {
	PositionRefreshDuration           *time.Duration `envconfig:"PROVIDER_POSITION_REFRESH_DURATION"`
	IsPositionRefreshDurationFloating bool           `envconfig:"PROVIDER_POSITION_REFRESH_DURATION_FLOATING"`
	TradersRefreshDuration            *time.Duration `envconfig:"PROVIDER_TRADERS_REFRESH_DURATION"`
}

func (c *Provider) Prepare() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}

	if c.PositionRefreshDuration == nil {
		positionRefreshDuration := defaultPositionRefreshDuration
		c.PositionRefreshDuration = &positionRefreshDuration
	}

	if c.IsPositionRefreshDurationFloating {
		rand.Seed(time.Now().UnixNano())
	}

	if c.TradersRefreshDuration == nil {
		tradersRefreshDuration := defaultTradersRefreshDuration
		c.TradersRefreshDuration = &tradersRefreshDuration
	}

	return nil
}

func (c *Provider) PositionRefreshTimeDuration() time.Duration {
	if c.IsPositionRefreshDurationFloating {
		maxDuration := *c.PositionRefreshDuration
		minDuration := maxDuration / 2

		return time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
	}

	return *c.PositionRefreshDuration
}

func (c *Provider) TradersRefreshTimeDuration() time.Duration {
	return *c.TradersRefreshDuration
}
