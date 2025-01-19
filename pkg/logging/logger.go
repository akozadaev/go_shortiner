package logging

import (
	"context"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type contextKey = string

const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

var conf = &Config{
	Encoding:    "console",
	Level:       zapcore.InfoLevel,
	Development: true,
}

type Config struct {
	Encoding        string
	Level           zapcore.Level
	Development     bool
	InfoFilename    string
	InfoMaxSize     int
	InfoMaxBackups  int
	InfoMaxAge      int
	InfoCompress    bool
	ErrorFilename   string
	ErrorMaxSize    int
	ErrorMaxBackups int
	ErrorMaxAge     int
	ErrorCompress   bool
}

// SetConfig sets given logging configs for DefaultLogger's logger.
// Must set configs before calling DefaultLogger()
func SetConfig(c *Config) {
	conf = &Config{
		Encoding:        c.Encoding,
		Level:           c.Level,
		Development:     c.Development,
		InfoFilename:    c.InfoFilename,
		InfoMaxSize:     c.InfoMaxSize,
		InfoMaxBackups:  c.InfoMaxBackups,
		InfoMaxAge:      c.InfoMaxAge,
		InfoCompress:    c.InfoCompress,
		ErrorFilename:   c.ErrorFilename,
		ErrorMaxSize:    c.ErrorMaxSize,
		ErrorMaxBackups: c.ErrorMaxBackups,
		ErrorMaxAge:     c.ErrorMaxAge,
		ErrorCompress:   c.ErrorCompress,
	}
}

func SetLevel(l zapcore.Level) {
	conf.Level = l
}

// NewLogger creates a new logger with the given log level
func NewLogger(conf *Config) *zap.SugaredLogger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	wsInfo := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.InfoFilename,
		MaxSize:    conf.InfoMaxSize, //MB
		MaxBackups: conf.InfoMaxBackups,
		MaxAge:     conf.InfoMaxAge, //days
		Compress:   conf.InfoCompress,
	})

	wsError := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.ErrorFilename,
		MaxSize:    conf.ErrorMaxSize, //MB
		MaxBackups: conf.ErrorMaxBackups,
		MaxAge:     conf.ErrorMaxAge, //days
		Compress:   conf.ErrorCompress,
	})
	coreInfo := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), wsInfo),
		lowPriority,
	)

	coreError := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), wsError),
		highPriority,
	)

	cores := zapcore.NewTee(coreInfo, coreError)

	logger := zap.New(cores)
	defer logger.Sync()

	return logger.Sugar()
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger(conf)
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	if gCtx, ok := ctx.(*gin.Context); ok {
		ctx = gCtx.Request.Context()
	}
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return DefaultLogger()
	}
	if gCtx, ok := ctx.(*gin.Context); ok && gCtx != nil {
		ctx = gCtx.Request.Context()
	}
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}
