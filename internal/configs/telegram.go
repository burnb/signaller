package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Telegram struct {
	Token  string `envconfig:"TELEGRAM_TOKEN"`
	ChatId int64  `envconfig:"TELEGRAM_CHAT_ID"`
}

// Prepare variables to static configuration
func (c *Telegram) Prepare() error {
	return envconfig.Process("", c)
}
