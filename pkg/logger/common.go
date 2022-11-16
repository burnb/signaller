package logger

type Config interface {
	MinimalLevel() string
}

type ConfigTelegram interface {
	IsEnabled() bool
	Token() string
	ChatId() int64
}
