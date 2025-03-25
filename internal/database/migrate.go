package database

import (
	"go_shurtiner/internal/app/model"
	"go_shurtiner/pkg/logging"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Link{}, &model.User{})
	if err != nil {
		logging.DefaultLogger().Errorf("failed to migrate database: %v", err)
		return err
	}
	//P@$$w0rd
	user := &model.User{Name: "Алексей", MiddleName: "Сергеевич", LastName: "Козадаев", Password: "$2a$12$uXR.vgCffldZK3ryULgx6u8ld.sntTBJZgH4KPHt9fWEHU8X38zoW", Email: "akozadaev@inbox.ru"}
	db.FirstOrCreate(user)
	return nil
}
