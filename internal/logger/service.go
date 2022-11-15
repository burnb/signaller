package logger

import (
	"fmt"
	logBasic "log"
	"strings"

	zap2tg "github.com/alfonmga/zap2telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/burnb/signaller/internal/configs"
)

type Creator struct {
	defaultLogger *zap.Logger
	cores         []zapcore.Core
	entries       map[string]*zap.Logger
}

func NewCreator(cfg configs.Logger, tgCfg configs.Telegram) (*Creator, error) {
	logsLevel := new(zapcore.Level)
	err := logsLevel.Set(cfg.MinimalLevel)
	if err != nil {
		return nil, err
	}

	lCfg := zap.NewDevelopmentConfig()
	lCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	lCfg.Level = zap.NewAtomicLevelAt(*logsLevel)
	lCfg.OutputPaths = []string{"stdout"}

	defaultLogger, err := lCfg.Build()
	if err != nil {
		return nil, err
	}
	cores := []zapcore.Core{defaultLogger.Core()}

	if tgCfg.IsEnabled() {
		telegramCore, err :=
			zap2tg.NewTelegramCore(
				*tgCfg.Token,
				[]int64{*tgCfg.ChatId},
				zap2tg.WithLevel(zapcore.WarnLevel),
				zap2tg.WithNotificationOn(
					[]zapcore.Level{zapcore.WarnLevel, zap.ErrorLevel, zap.PanicLevel, zap.FatalLevel},
				),
				zap2tg.WithParseMode(tgbotapi.ModeMarkdownV2),
				zap2tg.WithFormatter(func(e zapcore.Entry, fields []zapcore.Field) string {
					msg := fmt.Sprintf("*%s*\n%s\n", strings.ToUpper(e.Message), e.LoggerName)

					var msgFields string
					for _, field := range fields {
						enc := zapcore.NewMapObjectEncoder()
						field.AddTo(enc)
						for k, v := range enc.Fields {
							msgFields += fmt.Sprintf(
								"%s: %s\n",
								strings.ToUpper(k),
								tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, fmt.Sprintf("%+v", v)),
							)
						}
					}
					if msgFields != "" {
						msg += fmt.Sprintf("\n%s", msgFields)
					}

					return msg
				}),
			)
		if err != nil {
			return nil, err
		}

		cores = append(cores, telegramCore)
	}

	return &Creator{
		defaultLogger: defaultLogger,
		entries:       make(map[string]*zap.Logger),
		cores:         cores,
	}, nil
}

func (s *Creator) Create(named string) *zap.Logger {
	log, ok := s.entries[named]
	if !ok {
		log = zap.New(zapcore.NewTee(s.cores...))
		log.Named(named)
		s.entries[named] = log
	}

	zap.ReplaceGlobals(log)

	return log
}

func (s *Creator) Shutdown() {
	for _, log := range s.entries {
		if logErr := log.Sync(); logErr != nil {
			logBasic.Printf("unable to sync logger %v", logErr)
		}
	}
}
