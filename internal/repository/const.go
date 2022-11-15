package repository

import (
	"time"
)

const (
	loggerName = "repository"

	defaultMaxConn         = 5
	defaultConnMaxLifetime = time.Minute * 5
)
