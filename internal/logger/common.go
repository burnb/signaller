package logger

import (
	"github.com/burnb/signaller/internal/configs"
)

type Config interface {
	IsDebug() bool
	GetMinimalLogLevel() string
	TelegramCfg() configs.Telegram
}
