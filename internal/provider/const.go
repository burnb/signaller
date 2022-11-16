package provider

import (
	"time"
)

const (
	loggerName = "ProviderService"

	defaultPositionRefreshTime = 10 * time.Second
	defaultTradersRefreshTime  = 24 * time.Hour
)
