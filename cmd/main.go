package main

import (
	logBasic "log"
	"os"
	osSignal "os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/burnb/signaller/internal/configs"
	"github.com/burnb/signaller/internal/grpc"
	"github.com/burnb/signaller/internal/metric"
	"github.com/burnb/signaller/internal/provider"
	"github.com/burnb/signaller/internal/proxy"
	"github.com/burnb/signaller/internal/repository"
	"github.com/burnb/signaller/pkg/exchange/binance"
	"github.com/burnb/signaller/pkg/logger"
)

const loggerName = "Signaller"

func main() {
	if err := godotenv.Load(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logBasic.Fatal(err)
		}
	}

	cfg := &configs.App{}
	if err := cfg.Prepare(); err != nil {
		logBasic.Fatal(err)
	}

	logCreator, err := logger.NewCreator(&cfg.Logger, &cfg.Telegram)
	if err != nil {
		logBasic.Fatal(err)
	}
	defer logCreator.Shutdown()

	log := logCreator.Create(loggerName)

	repo := repository.NewMysql(cfg.Db, log)
	if err = repo.Init(); err != nil {
		log.Panic("unable to init repository", zap.Error(err))
	}
	defer repo.Shutdown()

	grpcSrv := grpc.NewServer(cfg.GRPC, log)
	if err = grpcSrv.Init(); err != nil {
		log.Panic("unable to init grpc server", zap.Error(err))
	}

	proxySrv := proxy.New(cfg.Proxy, log)
	if err = proxySrv.Init(); err != nil {
		log.Panic("unable to init proxy service", zap.Error(err))
	}

	exClient := binance.NewClient(log, proxySrv)
	providerSrv := provider.NewService(log, exClient, repo, grpcSrv)
	if err = providerSrv.Init(); err != nil {
		log.Panic("unable to init provider", zap.Error(err))
	}

	metricSrv := metric.New(cfg.Metric, providerSrv, log)
	metricSrv.Init()

	log.Warn(logger.ColorGreen.Fill("started"))
	defer func() { log.Warn(logger.ColorRed.Fill("shutdown")) }()

	c := make(chan os.Signal)
	osSignal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
