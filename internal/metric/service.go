package metric

import (
	"net/http"
	"net/http/pprof"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/burnb/signaller/internal/configs"
)

type Service struct {
	cfg       configs.Metric
	log       *zap.Logger
	startedAt time.Time
}

func New(cfg configs.Metric, log *zap.Logger) *Service {
	return &Service{cfg: cfg, log: log.Named(loggerName), startedAt: time.Now()}
}

func (s *Service) Init() {
	if !s.cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		ginzap.Ginzap(s.log, time.RFC3339, false),
		ginzap.RecoveryWithZap(s.log, true),
	)
	router.GET(s.cfg.HttpPath(), s.httpHandlerMain)
	router.GET("/pprof/", gin.WrapF(pprof.Index))
	router.GET("/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	router.GET("/pprof/profile", gin.WrapF(pprof.Profile))
	router.GET("/pprof/symbol", gin.WrapF(pprof.Symbol))
	router.GET("/pprof/trace", gin.WrapF(pprof.Trace))
	router.GET("/pprof/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	router.GET("/pprof/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	router.GET("/pprof/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	router.GET("/pprof/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
	router.GET("/pprof/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	router.GET("/pprof/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))

	server := &http.Server{
		Addr:         s.cfg.Address(),
		Handler:      router,
		ReadTimeout:  defaultHttpReadTimeout,
		WriteTimeout: defaultHttpWriteTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			s.log.Error("unable to listen and serve metric http server", zap.Error(err))
		}
	}()

	s.log.Debug("initiated successfully", zap.String("address", s.cfg.Address()))
}

func (s *Service) httpHandlerMain(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		Main{
			Uptime: time.Since(s.startedAt).String(),
		},
	)
}
