package repository

import (
	"time"
)

const (
	loggerName = "Repository"

	defaultMaxConn         = 5
	defaultConnMaxLifetime = time.Minute * 5
)
