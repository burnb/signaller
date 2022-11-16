package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Telegram struct {
	TelegramToken  string `envconfig:"TELEGRAM_TOKEN"`
	TelegramChatId int64  `envconfig:"TELEGRAM_CHAT_ID"`
}

func (c *Telegram) IsEnabled() bool {
	return c.TelegramToken != "" && c.TelegramChatId != 0
}

func (c *Telegram) Token() string {
	return c.TelegramToken
}

func (c *Telegram) ChatId() int64 {
	return c.TelegramChatId
}

func (c *Telegram) Prepare() error {
	return envconfig.Process("", c)
}
