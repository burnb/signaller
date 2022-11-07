package main

import (
	logBasic "log"
	"os"
	osSignal "os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/burnb/signaller/internal/configs"
	"github.com/burnb/signaller/internal/grpc"
	"github.com/burnb/signaller/internal/logger"
	"github.com/burnb/signaller/internal/metric"
	"github.com/burnb/signaller/internal/provider"
	"github.com/burnb/signaller/internal/proxy"
	"github.com/burnb/signaller/internal/repository"
	"github.com/burnb/signaller/pkg/exchange/clients/binance"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logBasic.Fatal(err)
	}

	cfg := &configs.App{}
	if err := cfg.Prepare(); err != nil {
		logBasic.Fatal(err)
	}

	logCreator, err := logger.NewCreator(cfg)
	if err != nil {
		logBasic.Fatal(err)
	}
	log := logCreator.Create("app")
	defer func() {
		if logErr := log.Sync(); logErr != nil {
			logBasic.Printf("unable to sync logger %v", logErr)
		}
	}()

	db, dbErr := sqlx.Open("mysql", cfg.Db.GetDatabaseDSN())
	if dbErr != nil {
		log.Fatal("unable to open db", zap.Error(dbErr))
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	dbClient := db
	defer func() { _ = dbClient.Close() }()

	repo := repository.NewMysql(dbClient)

	grpcSrv := grpc.NewServer(cfg.GRPCAddress(), log)
	if err = grpcSrv.Init(); err != nil {
		log.Fatal("unable to init grpc server", zap.Error(err))
	}

	proxySrv := proxy.New(&cfg.Proxy, log)
	if err = proxySrv.Init(); err != nil {
		log.Fatal("unable to init proxy service", zap.Error(err))
	}

	exClient := binance.NewClient(log, proxySrv)

	providerSrv := provider.NewService(log, exClient, repo, grpcSrv)
	if err = providerSrv.InitAndServe(); err != nil {
		log.Fatal("unable to init provider", zap.Error(err))
	}

	metricSrv := metric.New(&cfg.Metric, log)
	metricSrv.Init()

	log.Info(logger.ColorGreen.Fill("started"))

	c := make(chan os.Signal)
	osSignal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Warn(logger.ColorRed.Fill("shutdown"))
}
