package metric

import (
	"time"
)

const (
	ServiceName = "Metric service"

	DefaultHttpReadTimeout  = 5 * time.Second
	DefaultHttpWriteTimeout = 10 * time.Second
)
