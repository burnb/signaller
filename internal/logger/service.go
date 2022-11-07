package logger

import (
	"fmt"
	"strings"

	"github.com/alfonmga/zap2telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Creator struct {
	defaultLogger *zap.Logger
	cores         []zapcore.Core
	entries       map[string]*zap.Logger
}

func NewCreator(cfg Config) (*Creator, error) {
	logsLevel := new(zapcore.Level)
	err := logsLevel.Set(cfg.GetMinimalLogLevel())
	if err != nil {
		return nil, err
	}

	lCfg := zap.NewDevelopmentConfig()
	lCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	lCfg.Level = zap.NewAtomicLevelAt(*logsLevel)
	lCfg.OutputPaths = []string{"stdout"}
	if cfg.IsDebug() {
		lCfg.Level.SetLevel(zap.DebugLevel)
	}

	defaultLogger, err := lCfg.Build()
	if err != nil {
		return nil, err
	}
	cores := []zapcore.Core{defaultLogger.Core()}

	if cfg.TelegramCfg().Token != "" && cfg.TelegramCfg().ChatId != 0 {
		telegramCore, err :=
			zap2telegram.NewTelegramCore(
				cfg.TelegramCfg().Token,
				[]int64{cfg.TelegramCfg().ChatId},
				zap2telegram.WithLevel(zapcore.WarnLevel), // send only Info and above logs to Telegram
				zap2telegram.WithNotificationOn(
					[]zapcore.Level{zapcore.WarnLevel, zap.ErrorLevel, zap.PanicLevel, zap.FatalLevel}, // enable message notification only this levels
				),
				zap2telegram.WithFormatter(func(e zapcore.Entry, fields []zapcore.Field) string {
					escapedLoggerName := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, e.LoggerName)
					escapedCaller := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, e.Caller.TrimmedPath())
					escapedMessage := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, e.Message)
					msg := fmt.Sprintf(
						"[%s] Logger: %s\nCaller: %s\nMessage: *%s*",
						strings.ToUpper(e.Level.String()),
						escapedLoggerName,
						escapedCaller,
						escapedMessage,
					)
					// add fields to the message
					msgFields := ""
					for _, field := range fields {
						enc := zapcore.NewMapObjectEncoder()
						field.AddTo(enc)
						for k, v := range enc.Fields {
							escapedK := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, k)
							escapedV := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, fmt.Sprintf("%+v", v))
							msgField := fmt.Sprintf("%s\\=`%s`", escapedK, escapedV)
							if msgFields != "" {
								msgFields += " " // add leading space if there are already fields
							}
							msgFields += msgField
						}
					}
					if msgFields != "" {
						msg += fmt.Sprintf("\n%s", msgFields)
					}
					if e.Stack != "" {
						msg += fmt.Sprintf("\nLogger stacktrace: `%s`", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, e.Stack))
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
