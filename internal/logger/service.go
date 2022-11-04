package logger

import (
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

	return &Creator{
		defaultLogger: defaultLogger,
		entries:       make(map[string]*zap.Logger),
	}, nil
}

func (s *Creator) Create(named string) *zap.Logger {
	log, ok := s.entries[named]
	if !ok {
		log = zap.New(zapcore.NewTee(s.defaultLogger.Core()))
		log.Named(named)
		s.entries[named] = log
	}

	zap.ReplaceGlobals(log)

	return log
}
