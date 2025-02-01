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

	user := &model.User{Name: "Алексей", MiddleName: "Сергеевич", LastName: "Козадаев", Password: "$2a$12$RwH.0DgK2A5aICN8ahK9oOTJmS5rWFDQ45qUx.y4ySEZ8la5Hdypq"}
	db.FirstOrCreate(user)
	return nil
}
