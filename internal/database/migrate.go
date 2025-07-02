package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/pkg/logging"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, dataSourceName string) error {
	err := db.AutoMigrate(&model.Link{}, &model.User{}, &model.Task{}, &model.JobQueue{})
	if err != nil {
		logging.DefaultLogger().Errorf("failed to auto migrate database: %v", err)
		return err
	}
	// P@$$w0rd
	user := &model.User{Name: "Алексей", MiddleName: "Сергеевич", LastName: "Козадаев", Password: "$2a$12$uXR.vgCffldZK3ryULgx6u8ld.sntTBJZgH4KPHt9fWEHU8X38zoW", Email: "akozadaev@inbox.ru"}
	db.FirstOrCreate(user)
	err = gooseMigrate(dataSourceName)
	if err != nil {
		logging.DefaultLogger().Errorf("failed to goose migrate database: %v", err)
		return err
	}
	return nil
}

func gooseMigrate(dataSourceName string) error {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		logging.DefaultLogger().Errorf("failed to open DB to goosee migrate: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if err = goose.Up(db, "migrations"); err != nil {
		logging.DefaultLogger().Errorf("failed to goosee migrate UP: %v", err)
		return err
	}
	return err
}
