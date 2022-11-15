package provider

import (
	"time"
)

const (
	loggerName = "LeaderboardServiceWorker"

	defaultPositionRefreshTime = 10 * time.Second
	defaultTradersRefreshTime  = 24 * time.Hour
)
