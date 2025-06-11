package main

import (
	"context"
	"errors"
	"fmt"
	"go_shurtiner/internal/app"
	"go_shurtiner/internal/app/authentication"
	"go_shurtiner/internal/app/repository"
	"go_shurtiner/internal/database"
	"go_shurtiner/internal/http/middleware"
	"go_shurtiner/internal/job"
	"go_shurtiner/internal/queue"
	queueSvc "go_shurtiner/internal/queue/service"
	reportSvc "go_shurtiner/internal/report/service"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var serverCmd = &cobra.Command{
	Use: "server:go_shurtiner",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func main() {
	if err := serverCmd.Execute(); err != nil {
		log.Printf("failed to execute command. err: %v", err)
		os.Exit(1)
	}
}

func runApplication() {
	serverConfig, err := config.Load()
	if err != nil {
		log.Error().Stack().Err(err)
	}

	loggerLevel := zapcore.Level(serverConfig.LoggingConfig.Level)
	if !serverConfig.LoggingConfig.Development {
		loggerLevel = zapcore.ErrorLevel
	}

	logging.SetConfig(&logging.Config{
		Encoding:        serverConfig.LoggingConfig.Encoding,
		Level:           loggerLevel,
		InfoFilename:    serverConfig.LoggingConfig.Info.Filename,
		InfoMaxSize:     serverConfig.LoggingConfig.Info.MaxSize,
		InfoMaxBackups:  serverConfig.LoggingConfig.Info.MaxBackups,
		InfoMaxAge:      serverConfig.LoggingConfig.Info.MaxAge,
		InfoCompress:    serverConfig.LoggingConfig.Info.Compress,
		ErrorFilename:   serverConfig.LoggingConfig.Error.Filename,
		ErrorMaxSize:    serverConfig.LoggingConfig.Error.MaxSize,
		ErrorMaxBackups: serverConfig.LoggingConfig.Error.MaxBackups,
		ErrorMaxAge:     serverConfig.LoggingConfig.Error.MaxAge,
		ErrorCompress:   serverConfig.LoggingConfig.Error.Compress,
	})

	defer func() {
		_ = logging.DefaultLogger().Sync()
	}()

	application := fx.New(
		fx.Supply(serverConfig),
		fx.Supply(logging.DefaultLogger().Desugar()),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.StopTimeout(serverConfig.ServerConfig.GracefulShutdown+time.Second),
		fx.Provide(
			// setup database
			database.NewDatabase,
			// server
			newServer,
			repository.NewShortenRepository,
			repository.NewUserRepository,
			app.NewHandler,
			authentication.NewBasicAuth,
			// task queue
			repository.NewQueueRepository,
			queueSvc.NewQueueService,
			fx.Annotate(
				repository.NewQueueRepository,
				fx.As(new(queueSvc.QueueRepository)),
				//fx.As(new(reportSvc.ReportRepository)),
			),
			fx.Annotate(
				queueSvc.NewQueueService,
				fx.As(new(queue.QueueService)),
			),
			queue.NewQueue,
			// report
			reportSvc.NewReportService,
			repository.NewReportRepository, // нужен для newQueue
			fx.Annotate(
				repository.NewReportRepository,
				fx.As(new(reportSvc.ReportRepository)),
			),
		),
		fx.Invoke(
			newQueue,
			app.RouteV1,
			app.RouteV2,
			func(r *gin.Engine) {},
		),
	)
	application.Run()
}

func newServer(lc fx.Lifecycle, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.TimeoutMiddleware(cfg.ServerConfig.WriteTimeout))
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RestfulParamsMiddleware())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Start to rest api server :%d", cfg.ServerConfig.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logging.DefaultLogger().Errorw("failed to close http server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Info("Stopped rest api server")
			return srv.Shutdown(ctx)
		},
	})
	return r
}

func newQueue(
	lc fx.Lifecycle, cfg *config.Config, svc *queue.Queue,
	queueSvc *queueSvc.QueueService,
	repository repository.ReportRepository,
) {
	svc.AddJob("create.report",
		job.NewDataJob(context.Background(), repository, queueSvc, cfg.PrepareDataConfig),
	)
	svc.AddJob("prepare.data",
		job.NewDataJob(context.Background(), repository, queueSvc, cfg.PrepareDataConfig),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Info("Start task queue")
			go func() {
				if err := svc.Run(ctx); err != nil {
					logging.DefaultLogger().Errorw("failed to run task queue", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Info("Stopped task queue")
			return svc.Shutdown()
		},
	})
}
