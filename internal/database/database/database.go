package database

import (
	"go_shurtiner/internal/shorten/model"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database with given config
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var (
		db     *gorm.DB
		err    error
		logger = NewLogger(time.Second, true, zapcore.Level(cfg.DBConfig.LogLevel))
	)

	for i := 0; i <= 50; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DBConfig.DataSourceName), &gorm.Config{Logger: logger, TranslateError: true})
		// Migrate the schema
		db.AutoMigrate(&model.Link{})
		if err == nil {
			break
		}
		logging.DefaultLogger().Warnf("failed to open database: %v", err)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}

	rawDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	rawDB.SetMaxOpenConns(cfg.DBConfig.Pool.MaxOpen)
	rawDB.SetMaxIdleConns(cfg.DBConfig.Pool.MaxIdle)
	rawDB.SetConnMaxLifetime(cfg.DBConfig.Pool.MaxLifetime)

	return db, nil
}
