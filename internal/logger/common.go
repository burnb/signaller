package logger

type Config interface {
	IsDebug() bool
	GetMinimalLogLevel() string
}
