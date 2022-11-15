package metric

import (
	"time"
)

const (
	loggerName = "MetricService"

	defaultHttpReadTimeout  = 5 * time.Second
	defaultHttpWriteTimeout = 10 * time.Second
)
